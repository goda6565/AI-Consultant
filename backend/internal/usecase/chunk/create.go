package chunk

import (
	"context"
	"fmt"
	"io"

	chunkEntity "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/entity"
	chunkService "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/service"
	chunkValue "github.com/goda6565/ai-consultant/backend/internal/domain/chunk/value"
	documentEntity "github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	documentRepository "github.com/goda6565/ai-consultant/backend/internal/domain/document/repository"
	documentValue "github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	llm "github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	logger "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/errors"
	storagePort "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
	transactionPorts "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

type CreateChunkInputPort interface {
	Execute(ctx context.Context, input CreateChunkUseCaseInput, logger logger.Logger) (*CreateChunkOutput, error)
}

type CreateChunkUseCaseInput struct {
	DocumentID string
}

type CreateChunkOutput struct {
	NumCreated int
}

type CreateChunkInteractor struct {
	vectorUnitOfWork   transactionPorts.VectorUnitOfWork
	documentRepository documentRepository.DocumentRepository
	pdfParser          *chunkService.PdfParser
	csvAnalyzer        *chunkService.CsvAnalyzer
	chunker            *chunkService.Chunker
	storagePort        storagePort.StoragePort
	llmClient          llm.LLMClient
}

func NewCreateChunkUseCase(vectorUnitOfWork transactionPorts.VectorUnitOfWork, documentRepository documentRepository.DocumentRepository, pdfParser *chunkService.PdfParser, csvAnalyzer *chunkService.CsvAnalyzer, chunker *chunkService.Chunker, storagePort storagePort.StoragePort, llmClient llm.LLMClient) CreateChunkInputPort {
	return &CreateChunkInteractor{
		vectorUnitOfWork:   vectorUnitOfWork,
		documentRepository: documentRepository,
		pdfParser:          pdfParser,
		csvAnalyzer:        csvAnalyzer,
		chunker:            chunker,
		storagePort:        storagePort,
		llmClient:          llmClient,
	}
}

const maxRetryCount = 3

func (i *CreateChunkInteractor) Execute(ctx context.Context, input CreateChunkUseCaseInput, logger logger.Logger) (result *CreateChunkOutput, err error) {
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

	// mark as sync start
	document.MarkAsSyncStart()

	currentRetryCount := document.GetRetryCount()
	switch {
	case currentRetryCount > maxRetryCount: // retry count is less than max retry count
		document.MarkAsSyncFailed()
		if updateErr := i.updateDocument(ctx, document); updateErr != nil {
			logger.Error("failed to update document status", "error", updateErr)
			return nil, fmt.Errorf("failed to update document status: %w", updateErr)
		}
		return nil, errors.NewUseCaseError(errors.InternalError, "failed to create chunks: max retry count reached")
	default:
		// increment retry count
		document.IncrementRetryCount()
		if updateErr := i.updateDocument(ctx, document); updateErr != nil {
			logger.Error("failed to update document status", "error", updateErr)
			return nil, fmt.Errorf("failed to update document status: %w", updateErr)
		}
	}

	// download document
	reader, err := i.storagePort.Download(ctx, document.GetStorageInfo())
	if err != nil {
		return nil, fmt.Errorf("failed to download document: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()

	// process document
	logger.Info("processing document", "document_id", document.GetID().Value())
	var text string
	switch document.GetDocumentType() {
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
	logger.Info("chunking start", "document_id", document.GetID().Value())
	chunkerOutput, err := i.chunker.Execute(ctx, chunkService.ChunkerInput{Text: text})
	if err != nil {
		return nil, fmt.Errorf("failed to chunk document: %w", err)
	}

	// create embeddings with batch processing
	logger.Info("creating embeddings number", "number", len(chunkerOutput.Chunks), "document_id", document.GetID().Value())
	contents := make([]string, len(chunkerOutput.Chunks))
	for i, chunk := range chunkerOutput.Chunks {
		contents[i] = chunk.Content
	}

	// Process embeddings in batches to respect Vertex AI limits (max 250 per batch)
	const maxBatchSize = 100
	var allEmbeddings [][]float32

	for start := 0; start < len(contents); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(contents) {
			end = len(contents)
		}

		batchContents := contents[start:end]
		embeddingInput := llm.GenerateEmbeddingBatchInput{
			Texts: batchContents,
			// TODO: select model by user setting
			Config: llm.EmbeddingConfig{Provider: llm.VertexAI, Model: llm.GeminiEmbedding001},
		}
		embeddingOutput, err := i.llmClient.GenerateEmbeddingBatch(ctx, embeddingInput)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embeddings for batch %d-%d: %w", start, end-1, err)
		}

		allEmbeddings = append(allEmbeddings, embeddingOutput.Embeddings...)
	}

	// create chunks
	chunks := make([]*chunkEntity.Chunk, len(chunkerOutput.Chunks))
	for i, chunk := range chunkerOutput.Chunks {
		embedding := allEmbeddings[i]
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
	err = i.vectorUnitOfWork.WithTx(ctx, func(ctx context.Context) error {
		for _, chunk := range chunks {
			err := i.vectorUnitOfWork.ChunkRepository(ctx).Create(ctx, chunk)
			if err != nil {
				return fmt.Errorf("failed to create chunk: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		logger.Error("failed to create chunks", "error", err)
		return nil, errors.NewUseCaseError(errors.InternalError, "failed to create chunks")
	}

	// mark as sync done
	document.MarkAsSyncDone()
	err = i.updateDocument(ctx, document)
	if err != nil {
		logger.Error("failed to update document status", "error", err)
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return &CreateChunkOutput{NumCreated: len(chunks)}, nil
}

func (i *CreateChunkInteractor) updateDocument(ctx context.Context, document *documentEntity.Document) error {
	numUpdated, err := i.documentRepository.Update(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	if numUpdated != 1 {
		return errors.NewUseCaseError(errors.InternalError, "failed to update document")
	}
	return nil
}
