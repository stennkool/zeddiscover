// Package sync orchestrates the model-sync operation: for every
// openai_compatible provider in the Zed config, it fetches the model list
// from the provider's /models endpoint and updates available_models.
package sync

import (
	"fmt"

	"github.com/stenn/zed-modeldiscovery/internal/config"
	"github.com/stenn/zed-modeldiscovery/internal/model"
	"github.com/stenn/zed-modeldiscovery/internal/provider"
)

// Result summarises a single sync run.
type Result struct {
	Synced  int      // number of providers successfully synced
	Total   int      // total number of models written
	Skipped []string // providers skipped with reason
}

// Runner executes the sync operation.
type Runner struct {
	Repo    *config.Repository
	Fetcher provider.Fetcher
}

// NewRunner creates a Runner with the given dependencies.
func NewRunner(repo *config.Repository, fetcher provider.Fetcher) *Runner {
	return &Runner{Repo: repo, Fetcher: fetcher}
}

// Run performs the full sync. If dryRun is true, providers are fetched but
// nothing is written to disk. If providerFilter is non-empty, only that
// provider is synced; others are left untouched.
func (r *Runner) Run(dryRun bool, providerFilter string) (*Result, error) {
	cfg, err := r.Repo.Read()
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	result := &Result{}
	providers, ok := r.extractProviders(cfg)
	if !ok {
		return result, nil
	}

	for name, raw := range providers {
		if providerFilter != "" && name != providerFilter {
			continue
		}

		provMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		apiURL, _ := provMap["api_url"].(string)
		if apiURL == "" {
			result.Skipped = append(result.Skipped,
				fmt.Sprintf("%s (no api_url)", name))
			continue
		}

		if dryRun {
			fmt.Printf("[dry-run] %s: fetching …\n", name)
		}

		models, err := r.Fetcher.FetchModels(apiURL)
		if err != nil {
			result.Skipped = append(result.Skipped,
				fmt.Sprintf("%s (%v)", name, err))
			continue
		}

		zedModels := r.filterAndMap(models)
		provMap["available_models"] = zedModels
		result.Synced++
		result.Total += len(zedModels)

		if dryRun {
			fmt.Printf("  → %d text models\n", len(zedModels))
		}
	}

	if !dryRun && result.Synced > 0 {
		if err := r.Repo.Write(cfg); err != nil {
			return nil, fmt.Errorf("write config: %w", err)
		}
	}

	return result, nil
}

// extractProviders navigates cfg → language_models → openai_compatible.
func (r *Runner) extractProviders(cfg config.ZedConfig) (map[string]any, bool) {
	lmRaw, ok := cfg["language_models"]
	if !ok {
		return nil, false
	}
	lm, ok := lmRaw.(map[string]any)
	if !ok {
		return nil, false
	}
	ocRaw, ok := lm["openai_compatible"]
	if !ok {
		return nil, false
	}
	providers, ok := ocRaw.(map[string]any)
	return providers, ok
}

// filterAndMap keeps only text-output models and converts them to the
// Zed available_model format.
func (r *Runner) filterAndMap(apiModels []model.APIModel) []model.AvailableModel {
	var out []model.AvailableModel
	for _, m := range apiModels {
		if !m.IsTextOnly() {
			continue
		}
		out = append(out, model.MapToAvailableModel(m))
	}
	return out
}
