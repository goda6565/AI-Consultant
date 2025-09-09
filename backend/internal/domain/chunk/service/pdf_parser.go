package service

import (
	"context"
	"io"

	"github.com/goda6565/ai-consultant/backend/internal/domain/ocr"
)

type PdfParserInput struct {
	Reader io.ReadCloser
}

type PdfParserOutput struct {
	Text string
}

type PdfParser struct {
	ocrClient ocr.OcrClient
}

func NewPdfParserService(ocrClient ocr.OcrClient) *PdfParser {
	return &PdfParser{ocrClient: ocrClient}
}

func (dp *PdfParser) Execute(ctx context.Context, input PdfParserInput) (*PdfParserOutput, error) {
	ocrInput := ocr.OcrInput{
		Extension: ocr.OCRDocumentExtensionPDF,
		Reader:    input.Reader,
	}
	ocrOutput, err := dp.ocrClient.ExtractText(ctx, ocrInput)
	if err != nil {
		return nil, err
	}
	return &PdfParserOutput{Text: ocrOutput.ExtractedText}, nil
}
