package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ParseJSONUnmarshalError converts JSON unmarshal errors into more user-friendly error messages.
func ParseJSONUnmarshalError(err error) error {
	// Check if it's a JSON unmarshal error
	var unmarshalErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	switch {
	case errors.As(err, &unmarshalErr):
		return fmt.Errorf("invalid value for field %q: expected %q but got %q",
			unmarshalErr.Field, unmarshalErr.Type.String(), unmarshalErr.Value)
	case errors.As(err, &syntaxErr):
		return fmt.Errorf("invalid input provided to the tool")
	default:
		return err
	}
}
