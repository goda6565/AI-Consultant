package ocr

import (
	"context"
	"fmt"
	"io"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"

	"github.com/goda6565/ai-consultant/backend/internal/domain/ocr"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
)

type DocumentAIClient struct {
	client *documentai.DocumentProcessorClient
	env    *environment.Environment
}

func NewDocumentAIClient(ctx context.Context, env *environment.Environment) *DocumentAIClient {
	client, err := documentai.NewDocumentProcessorClient(ctx)
	if err != nil {
		panic(err)
	}
	return &DocumentAIClient{client: client, env: env}
}

func (c *DocumentAIClient) ExtractText(ctx context.Context, input ocr.OcrInput) (*ocr.OcrOutput, error) {
	projectID := c.env.ProjectID
	location := c.env.DocumentAILocation
	processorID := c.env.ProcessorID
	data, err := io.ReadAll(input.Reader)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to read data: %v", err))
	}
	var mimeType string
	switch input.Extension {
	case ocr.OCRDocumentExtensionPDF:
		mimeType = "application/pdf"
	case ocr.OCRDocumentExtensionImage:
		mimeType = "image/png"
	}
	request := &documentaipb.ProcessRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", projectID, location, processorID),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: mimeType,
			},
		},
	}

	resp, err := c.client.ProcessDocument(ctx, request)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to process document: %v", err))
	}

	document := resp.GetDocument()
	return &ocr.OcrOutput{ExtractedText: document.GetText()}, nil
}
