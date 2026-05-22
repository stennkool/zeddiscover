// Package config handles reading, writing, and backing up the Zed editor
// settings.json file. It understands JSONC (trailing commas) and preserves
// the rest of the document untouched.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ZedConfig is a loosely-typed representation of the full settings.json.
// Only the language_models.openai_compatible subtree is ever modified;
// everything else is carried through transparently.
type ZedConfig map[string]any

// Repository provides access to the persisted settings.json.
type Repository struct {
	Path string
}

// NewRepository creates a Repository for the given config path.
func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

// DefaultPath returns the standard Zed config location.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "zed", "settings.json")
}

// Read loads and parses the settings file, handling JSONC trailing commas.
func (r *Repository) Read() (ZedConfig, error) {
	raw, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", r.Path, err)
	}
	cleaned := cleanJSONC(raw)

	var cfg ZedConfig
	if err := json.Unmarshal(cleaned, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", r.Path, err)
	}
	return cfg, nil
}

// Write persists the config and writes a .bak backup first.
func (r *Repository) Write(cfg ZedConfig) error {
	if err := r.backup(); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	// trailing newline to match typical JSON files
	out := append(bytes.TrimRight(buf.Bytes(), "\n"), '\n')
	return os.WriteFile(r.Path, out, 0o644)
}

// Backup copies the current file to <path>.bak.
func (r *Repository) backup() error {
	src, err := os.ReadFile(r.Path)
	if err != nil {
		return err // file doesn't exist yet — nothing to back up
	}
	return os.WriteFile(r.Path+".bak", src, 0o644)
}

// ─── JSONC pre-processing ───────────────────────────────────────────────────

var reTrailingComma = regexp.MustCompile(`,(\s*[}\]])`)

// cleanJSONC removes JSONC trailing commas so Go's json decoder can parse.
func cleanJSONC(raw []byte) []byte {
	return reTrailingComma.ReplaceAll(raw, []byte("$1"))
}
