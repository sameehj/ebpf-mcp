package types

import (
	"bytes"
	"encoding/json"
)

// StrictUnmarshal applies strict decoding to any JSON-compatible input into an arbitrary struct.
// Useful when the target is dynamically passed.
func StrictUnmarshal(input any, out any) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

// StrictUnmarshalTyped is the generic version of StrictUnmarshal.
// It returns a typed pointer to the decoded struct, or an error.
func StrictUnmarshalTyped[T any](input any) (*T, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	var out T
	if err := decoder.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
