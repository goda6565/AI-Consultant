package scraper

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

const MaxScrapeContentLength = 5000

type ScraperClient struct {
	httpClient *http.Client
}

func NewScraperClient() *ScraperClient {
	return &ScraperClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *ScraperClient) Scrape(ctx context.Context, url string) (string, error) {
	logger := logger.GetLogger(ctx)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d for %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	text := ""
	doc.Find("h1, h2, h3, p").Each(func(i int, s *goquery.Selection) {
		text += s.Text() + "\n"
	})

	length := utf8.RuneCountInString(text)
	if length > MaxScrapeContentLength {
		runes := []rune(text)
		text = string(runes[:MaxScrapeContentLength])
	}
	return text, nil
}
