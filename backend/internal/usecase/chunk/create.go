package chunk

import (
	"context"
	"fmt"
	"io"

	chunkEntity "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/entity"
	chunkRepository "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/repository"
	chunkService "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/service"
	chunkValue "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/value"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	documentValue "github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	llm "github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	errors "github.com/goda6565/ai-consultant/backend/internal/usecase/error"
	storagePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
)

type CreateChunkInputPort interface {
	Execute(ctx context.Context, input CreateChunkUseCaseInput) (*CreateChunkOutput, error)
}

type CreateChunkUseCaseInput struct {
	DocumentID string
}

type CreateChunkOutput struct {
	NumCreated int
}

type CreateChunkInteractor struct {
	chunkRepository    chunkRepository.ChunkRepository
	documentRepository documentRepository.DocumentRepository
	pdfParser          chunkService.PdfParser
	csvAnalyzer        chunkService.CsvAnalyzer
	chunker            chunkService.Chunker
	storagePort        storagePort.StoragePort
	llmClient          llm.LLMClient
}

func NewCreateChunkUseCase(chunkRepository chunkRepository.ChunkRepository, documentRepository documentRepository.DocumentRepository, pdfParser chunkService.PdfParser, csvAnalyzer chunkService.CsvAnalyzer, chunker chunkService.Chunker) CreateChunkInputPort {
	return &CreateChunkInteractor{chunkRepository: chunkRepository, documentRepository: documentRepository, pdfParser: pdfParser, csvAnalyzer: csvAnalyzer, chunker: chunker}
}

func (i *CreateChunkInteractor) Execute(ctx context.Context, input CreateChunkUseCaseInput) (*CreateChunkOutput, error) {
	// find document
	documentID, err := sharedValue.NewID(input.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("failed to create document id: %w", err)
	}
	document, err := i.documentRepository.FindById(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if document == nil {
		return nil, errors.NewUseCaseError(errors.NotFoundError, "document not found")
	}

	// download document
	reader, err := i.storagePort.Download(ctx, document.GetStoragePath())
	if err != nil {
		return nil, fmt.Errorf("failed to download document: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()

	// process document
	var text string
	switch document.GetDocumentExtension() {
	case documentValue.DocumentExtensionPDF:
		// parsed by ocr
		pdfParserOutput, err := i.pdfParser.Execute(ctx, chunkService.PdfParserInput{Reader: reader})
		if err != nil {
			return nil, fmt.Errorf("failed to parse document: %w", err)
		}
		text = pdfParserOutput.Text
	case documentValue.DocumentExtensionMarkdown:
		b, err := io.ReadAll(reader)
		if err != nil {
			return nil, errors.NewUseCaseError(errors.InternalError, "failed to read document")
		}
		text = string(b)
	case documentValue.DocumentExtensionCSV:
		// generate summary by llm
		csvAnalyzerOutput, err := i.csvAnalyzer.Execute(ctx, chunkService.CsvAnalyzerInput{Reader: reader})
		if err != nil {
			return nil, errors.NewUseCaseError(errors.InternalError, "failed to summarize csv")
		}
		text = csvAnalyzerOutput.Text
	default:
		return nil, errors.NewUseCaseError(errors.InternalError, "invalid document extension")
	}

	// chunk document
	chunkerOutput, err := i.chunker.Execute(ctx, chunkService.ChunkerInput{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to chunk document: %w", err)
	}

	// create embeddings
	contents := make([]string, len(chunkerOutput.Chunks))
	for i, chunk := range chunkerOutput.Chunks {
		contents[i] = chunk.Content
	}
	embeddingInput := llm.GenerateEmbeddingBatchInput{
		Texts:  contents,
		Config: llm.EmbeddingConfig{Provider: llm.OpenAI, Model: llm.EmbeddingModelOpenAIEmbeddings},
	}
	embeddingOutput, err := i.llmClient.GenerateEmbeddingBatch(ctx, embeddingInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// create chunks
	chunks := make([]*chunkEntity.Chunk, len(chunkerOutput.Chunks))
	for i, chunk := range chunkerOutput.Chunks {
		embedding := embeddingOutput.Embeddings[i]
		id, err := sharedValue.NewID(uuid.NewUUID())
		if err != nil {
			return nil, fmt.Errorf("failed to create id: %w", err)
		}
		content, err := chunkValue.NewContent(chunk.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to create content: %w", err)
		}
		parentContent, err := chunkValue.NewContent(chunk.ParentContent)
		if err != nil {
			return nil, fmt.Errorf("failed to create parent content: %w", err)
		}
		embeddingValue, err := chunkValue.NewEmbedding(embedding)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding: %w", err)
		}
		chunks[i] = chunkEntity.NewChunk(id, document.GetID(), content, parentContent, embeddingValue)
	}

	// create chunks
	for _, chunk := range chunks {
		err := i.chunkRepository.Create(ctx, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to create chunk: %w", err)
		}
	}

	return &CreateChunkOutput{NumCreated: 1}, nil
}
