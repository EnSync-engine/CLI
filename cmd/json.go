package cmd

import (
	"encoding/json"
	"io"
)

// printJSON writes the value as indented JSON to the writer.
func printJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
