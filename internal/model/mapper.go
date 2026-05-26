package model

import "slices"

// ─── Filtering predicates ───────────────────────────────────────────────────

// IsTextOnly returns true when the model only produces text output.
// Models that also output images, audio, or video are excluded.
func (m APIModel) IsTextOnly() bool {
	if len(m.outputMods()) == 0 {
		return true // no modalities listed → assume text-only
	}
	for _, mod := range m.outputMods() {
		if mod != "text" {
			return false
		}
	}
	return true
}

// ─── Capability detection ───────────────────────────────────────────────────

func hasImageInput(m APIModel) bool {
	return slices.Contains(m.inputMods(), "image")
}

func hasToolSupport(m APIModel) bool {
	return slices.Contains(m.supportedParams(), "tools")
}

func hasReasoningSupport(m APIModel) bool {
	return slices.Contains(m.supportedParams(), "reasoning") ||
		slices.Contains(m.supportedParams(), "include_reasoning")
}

// ─── Mapping ────────────────────────────────────────────────────────────────

// MapToAvailableModel converts an APIModel into the Zed AvailableModel format.
func MapToAvailableModel(m APIModel) AvailableModel {
	maxTokens := m.ContextLength
	maxCompletion := maxTokens // sensible fallback when provider doesn't specify
	if m.TopProvider.MaxCompletionTokens != nil {
		maxCompletion = *m.TopProvider.MaxCompletionTokens
	} else if m.MaxOutputTokens > 0 {
		maxCompletion = m.MaxOutputTokens
	}

	return AvailableModel{
		Name:                m.ID,
		MaxTokens:           maxTokens,
		MaxOutputTokens:     maxCompletion,
		MaxCompletionTokens: maxCompletion,
		Capabilities: Capabilities{
			Tools:                hasToolSupport(m),
			Images:               hasImageInput(m),
			ParallelToolCalls:    hasToolSupport(m),
			PromptCacheKey:       false,
			ChatCompletions:      true,
			InterleavedReasoning: hasReasoningSupport(m),
		},
	}
}
