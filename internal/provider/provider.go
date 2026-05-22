package provider

import "github.com/stenn/zed-modeldiscovery/internal/model"

// Fetcher retrieves available models from a provider's /models endpoint.
// Different providers may use different response formats; each concrete
// implementation handles its own wire format and returns the common
// model.APIModel slice.
type Fetcher interface {
	// FetchModels calls the provider's API and returns the parsed model list.
	// The apiURL is the base URL as configured in Zed's settings.json
	// (e.g. "https://api.kilo.ai/api/gateway").
	FetchModels(apiURL string) ([]model.APIModel, error)
}
