package ocr

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/goda6565/ai-consultant/backend/internal/domain/ocr"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
)

func TestDocumentAIClient_ExtractText(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	location := os.Getenv("DOCUMENT_AI_LOCATION")
	processorID := os.Getenv("DOCUMENT_AI_PROCESSOR_ID")
	if projectID == "" || location == "" || processorID == "" {
		t.Skip("GOOGLE_CLOUD_PROJECT_ID, DOCUMENT_AI_LOCATION, DOCUMENT_AI_PROCESSOR_ID are not set")
	}
	client := NewDocumentAIClient(context.Background(), &environment.Environment{
		GoogleCloudEnvironment: environment.GoogleCloudEnvironment{
			ProjectID: projectID,
		},
		DocumentAIEnvironment: environment.DocumentAIEnvironment{
			DocumentAILocation: location,
			ProcessorID:        processorID,
		},
	})
	// URLからPDFをダウンロード
	resp, err := http.Get("https://www.pref.aichi.jp/kenmin/shohiseikatsu/education/pdf/student_guide.pdf")
	if err != nil {
		t.Fatalf("failed to download PDF: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to download PDF: status code %d", resp.StatusCode)
	}

	// レスポンスボディをio.Readerとして使用
	reader := resp.Body
	input := ocr.OcrInput{
		Extension: ocr.OCRDocumentExtensionPDF,
		Reader:    reader,
	}
	ocrOutput, err := client.ExtractText(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to extract text: %v", err)
	}
	t.Logf("ocrOutput: %v", ocrOutput)
}
