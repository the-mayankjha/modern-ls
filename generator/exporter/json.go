// Package exporter writes the parsed icon database to assets/icons.json and
// can also read it back for code generation.
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/the-mayankjha/modern-ls/generator/parser"
)

// IconDB is the full icon database written to assets/icons.json.
// All maps are keyed by the lookup string (filename or extension).
type IconDB struct {
	// SchemaVersion allows future migrations without breaking old readers.
	SchemaVersion int `json:"schema_version"`
	// ByFilename maps exact filenames (e.g. "go.mod", "Dockerfile") to icons.
	ByFilename map[string]parser.Icon `json:"by_filename"`
	// ByExtension maps file extensions without the leading dot (e.g. "go", "ts") to icons.
	ByExtension map[string]parser.Icon `json:"by_extension"`
	// ByFolder maps specific folder names to icons.
	ByFolder map[string]parser.Icon `json:"by_folder"`
	// Defaults holds fallback icons for generic kinds: "file", "dir", "diropen",
	// "hiddenfile", "hiddendir", "exe".
	Defaults map[string]parser.Icon `json:"defaults"`
}

// WriteJSON encodes db as indented JSON to w.
// HTML escaping is disabled so that Nerd Font glyphs survive the round-trip
// unmodified (json.Encoder would otherwise escape characters above U+007F).
func WriteJSON(w io.Writer, db IconDB) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(db); err != nil {
		return fmt.Errorf("exporter: encode: %w", err)
	}
	return nil
}

// ReadJSON decodes an IconDB from r.
func ReadJSON(r io.Reader) (IconDB, error) {
	var db IconDB
	dec := json.NewDecoder(r)
	if err := dec.Decode(&db); err != nil {
		return IconDB{}, fmt.Errorf("exporter: decode: %w", err)
	}
	return db, nil
}

// SortedKeys returns the sorted keys of a map for deterministic code generation.
// Determinism matters so that the generated file does not produce spurious diffs.
func SortedKeys(m map[string]parser.Icon) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
