package googlesearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	searchClient "github.com/goda6565/ai-consultant/backend/internal/domain/search"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type GoogleSearchClient struct {
	env *environment.Environment
}

func NewGoogleSearchClient(env *environment.Environment) searchClient.WebSearchClient {
	return &GoogleSearchClient{env: env}
}

type googleSearchItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type googleSearchResponse struct {
	Items []googleSearchItem `json:"items"`
}

func (c *GoogleSearchClient) Search(ctx context.Context, input searchClient.WebSearchInput) (*searchClient.WebSearchOutput, error) {
	logger := logger.GetLogger(ctx)
	// create params
	params := url.Values{}
	params.Set("key", c.env.CustomSearchAPIKey)
	params.Set("cx", c.env.SearchEngineID)
	params.Set("q", input.Query)

	// HTML only
	params.Set("fileType", "html")
	params.Set("hq", "filetype:html")

	endpoint := fmt.Sprintf("%s?%s", c.env.SearchEndpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create request: %v", err))
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("google search request failed: %v", err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("google search request failed: %v", resp.StatusCode))
	}

	var response googleSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to decode response: %v", err))
	}

	var htmlItems []googleSearchItem
	for _, item := range response.Items {
		if strings.HasSuffix(strings.ToLower(item.Link), ".html") {
			htmlItems = append(htmlItems, item)
		}
	}

	var returnResponses []googleSearchItem
	source := htmlItems
	if len(source) == 0 {
		source = response.Items
	}
	if input.MaxNumResults > 0 && input.MaxNumResults < len(source) {
		returnResponses = source[:input.MaxNumResults]
	} else {
		returnResponses = source
	}

	results := []searchClient.WebSearchResult{}
	for _, item := range returnResponses {
		results = append(results, searchClient.WebSearchResult{
			Title:   item.Title,
			Snippet: item.Snippet,
			URL:     item.Link,
		})
	}

	return &searchClient.WebSearchOutput{Results: results}, nil
}
