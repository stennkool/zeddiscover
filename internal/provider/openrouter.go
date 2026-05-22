package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stenn/zed-modeldiscovery/internal/model"
)

// OpenRouterFetcher handles the OpenRouter / Kilo / Synthetic style
// /models endpoint. Responses look like:
//
//	{"data": [{"id": "...", "context_length": ..., ...}, ...]}
type OpenRouterFetcher struct {
	client *http.Client
}

// NewOpenRouterFetcher creates a fetcher with a sensible default HTTP client.
func NewOpenRouterFetcher() *OpenRouterFetcher {
	return &OpenRouterFetcher{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchModels implements the Fetcher interface. It appends "/models" to the
// provider's apiURL and parses the standard OpenRouter-style response.
func (f *OpenRouterFetcher) FetchModels(apiURL string) ([]model.APIModel, error) {
	url := strings.TrimRight(apiURL, "/") + "/models"

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("http %d from %s: %s", resp.StatusCode, url,
			strings.TrimSpace(string(body)))
	}

	var envelope model.ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode %s: %w", url, err)
	}
	return envelope.Data, nil
}
