package luno

import (
	"errors"
	"fmt"
)

// Error is a Luno API error.
type Error struct {
	// Code can be used to identify errors even if the error message is
	// localised.
	Code string `json:"error_code"`

	// Message may be localised for authenticated API calls.
	Message string `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

func (e Error) ErrCode() string {
	return e.Code
}

// IsErrorCode returns whether an error is identifiable by a given code. This can be used to handle luno.Client errors.
// Any other errors will cause this to return false.
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	var lErr Error
	if errors.As(err, &lErr) {
		return lErr.ErrCode() == code
	}

	return false
}
