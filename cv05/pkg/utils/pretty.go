package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PrettyPrintJSON pretty prints the given value as JSON.
func PrettyPrintJSON(v any) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}
