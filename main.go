// sync_zed_models keeps the Zed editor's language model definitions
// in sync with live provider /models endpoints.
//
// For every openai_compatible provider in ~/.config/zed/settings.json
// that has an api_url, it fetches <api_url>/models and updates the
// provider's available_models list.  Only text-output models are kept;
// image/audio/video generators are filtered out.
//
// Usage:
//
//	go run .                        # sync all providers
//	go run . --dry-run              # preview only
//	go run . --provider kilo        # sync only the "kilo" provider
//	go run . --config <path>        # custom config location
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stenn/zed-modeldiscovery/internal/config"
	"github.com/stenn/zed-modeldiscovery/internal/provider"
	"github.com/stenn/zed-modeldiscovery/internal/sync"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "preview changes without writing")
	customPath := flag.String("config", "", "path to settings.json (default: ~/.config/zed/settings.json)")
	providerName := flag.String("provider", "", "sync only this provider (default: all)")
	flag.Parse()

	cfgPath := *customPath
	if cfgPath == "" {
		cfgPath = config.DefaultPath()
	}

	repo := config.NewRepository(cfgPath)
	fetcher := provider.NewOpenRouterFetcher()
	runner := sync.NewRunner(repo, fetcher)

	result, err := runner.Run(*dryRun, *providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Printf("\n[dry-run] summary: %d provider(s) synced, %d total models\n",
			result.Synced, result.Total)
	} else {
		fmt.Printf("Synced %d provider(s) (%d total models) in %s\n",
			result.Synced, result.Total, cfgPath)
	}

	if len(result.Skipped) > 0 {
		fmt.Println("\nSkipped:")
		for _, s := range result.Skipped {
			fmt.Printf("  • %s\n", s)
		}
	}
}
