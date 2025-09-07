package ocr

import (
	"context"
	"io"
)

type OCRDocumentExtension string

const (
	OCRDocumentExtensionPDF   OCRDocumentExtension = "pdf"
	OCRDocumentExtensionImage OCRDocumentExtension = "png"
)

type OcrInput struct {
	Extension OCRDocumentExtension
	Reader    io.Reader
}

type OcrOutput struct {
	ExtractedText string
}

type OcrClient interface {
	ExtractText(ctx context.Context, input OcrInput) (*OcrOutput, error)
}
