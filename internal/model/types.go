// Package model defines the domain types used across the application:
// the Zed editor's available_model representation, and the API model
// representation returned by provider /models endpoints.
package model

// Capabilities describes what a model can do.
type Capabilities struct {
	Tools                bool `json:"tools"`
	Images               bool `json:"images"`
	ParallelToolCalls    bool `json:"parallel_tool_calls"`
	PromptCacheKey       bool `json:"prompt_cache_key"`
	ChatCompletions      bool `json:"chat_completions"`
	InterleavedReasoning bool `json:"interleaved_reasoning"`
}

// AvailableModel is the Zed editor's representation of a single model
// inside a provider's available_models list.
type AvailableModel struct {
	Name                string       `json:"name"`
	MaxTokens           int64        `json:"max_tokens"`
	MaxOutputTokens     int64        `json:"max_output_tokens,omitempty"`
	MaxCompletionTokens int64        `json:"max_completion_tokens,omitempty"`
	Capabilities        Capabilities `json:"capabilities"`
}

// APIModel is a single model as returned by an OpenRouter-style /models endpoint.
type APIModel struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	ContextLength   int64        `json:"context_length"`
	Architecture    Architecture `json:"architecture"`
	TopProvider     TopProvider  `json:"top_provider"`
	SupportedParams []string     `json:"supported_parameters"`
}

// Architecture describes input/output modalities.
type Architecture struct {
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
}

// TopProvider holds provider-level details like max completion tokens.
type TopProvider struct {
	MaxCompletionTokens *int64 `json:"max_completion_tokens"`
}

// ModelsResponse is the top-level envelope from a /models endpoint.
type ModelsResponse struct {
	Data []APIModel `json:"data"`
}
