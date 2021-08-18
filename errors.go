package luno

import (
	"fmt"
)

// lunoError is a Luno API error.
type lunoError struct {
	// ErrorCode can be used to identify errors even if the error message is
	// localised.
	ErrorCode string `json:"error_code"`

	// Message may be localised for authenticated API calls.
	Message string `json:"error"`
}

func (e lunoError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.ErrorCode)
}

func (e lunoError) Code() string {
	return e.ErrorCode
}

// IsErrorCode returns whether an error is identifiable by a given code. This can be used to handle luno.Client errors.
// Any other errors will cause this to return false.
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	if lErr, ok := err.(lunoError); ok {
		return lErr.Code() == code
	}

	return false
}
