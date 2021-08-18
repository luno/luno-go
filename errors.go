package luno

import (
	"fmt"
)

// Error shouldn't be used anymore. Please use IsErrorCode if you want to check for a specific error code.
// Deprecated
type Error struct {
	lunoError
}

// lunoError is a Luno API error.
type lunoError struct {
	// Code can be used to identify errors even if the error message is
	// localised.
	Code string `json:"error_code"`

	// Message may be localised for authenticated API calls.
	Message string `json:"error"`
}

func (e lunoError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

func (e lunoError) ErrCode() string {
	return e.Code
}

// IsErrorCode returns whether an error is identifiable by a given code. This can be used to handle luno.Client errors.
// Any other errors will cause this to return false.
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	if lErr, ok := err.(lunoError); ok {
		return lErr.ErrCode() == code
	}

	return false
}
