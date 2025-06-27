// pkg/types/unmarshal.go
package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// StrictUnmarshal safely unmarshals input with proper nil handling
func StrictUnmarshal(input map[string]interface{}, target interface{}) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	}

	if target == nil {
		return fmt.Errorf("target is nil")
	}

	// Check if target is a pointer
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	// Check if the pointer points to nil
	if rv.IsNil() {
		return fmt.Errorf("target pointer is nil")
	}

	// Use JSON marshaling/unmarshaling for safe conversion
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to target: %w", err)
	}

	return nil
}

// SafeStringFromInterface extracts a string from an interface{} with nil checks
func SafeStringFromInterface(value interface{}, fieldName string) (string, error) {
	if value == nil {
		return "", fmt.Errorf("%s is nil", fieldName)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string, got %T", fieldName, value)
	}

	return str, nil
}

// SafeMapFromInterface extracts a map from an interface{} with nil checks
func SafeMapFromInterface(value interface{}, fieldName string) (map[string]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("%s is nil", fieldName)
	}

	m, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be a map, got %T", fieldName, value)
	}

	return m, nil
}

// SafeBoolFromInterface extracts a bool from an interface{} with nil checks
func SafeBoolFromInterface(value interface{}, fieldName string) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("%s is nil", fieldName)
	}

	b, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean, got %T", fieldName, value)
	}

	return b, nil
}

// SafeFloat64FromInterface extracts a float64 from an interface{} with nil checks
func SafeFloat64FromInterface(value interface{}, fieldName string) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("%s is nil", fieldName)
	}

	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("%s must be a number, got %T", fieldName, value)
	}

	return f, nil
}

// SafeSliceFromInterface extracts a slice from an interface{} with nil checks
func SafeSliceFromInterface(value interface{}, fieldName string) ([]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("%s is nil", fieldName)
	}

	s, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an array, got %T", fieldName, value)
	}

	return s, nil
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
