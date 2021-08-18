package luno

import (
	"fmt"
	"testing"
)

func TestLunoError(t *testing.T) {
	err := lunoError{
		ErrorCode: "ErrFoo",
		Message:   "Example error message!",
	}

	 if err.Error() != "Example error message! (ErrFoo)" {
	 	t.Errorf("Unexpected error formatting")
	 }

	 if err.Code() != "ErrFoo" {
	 	t.Errorf("Unexpected error code")
	 }
}

func TestIsErrorCode(t *testing.T) {
	testCases := []struct {
		name string
		err error
		code string
		exp bool
	}{
		{
			name: "nil error",
			err:  nil,
			code: "ErrFoo",
			exp:  false,
		},
		{
			name: "luno error with code present",
			err:  lunoError{
				ErrorCode: "ErrFoo",
				Message:   "example error message",
			},
			code: "ErrFoo",
			exp:  true,
		},
		{
			name: "luno error with different code",
			err:  lunoError{
				ErrorCode: "ErrBar",
				Message:   "example error message",
			},
			code: "ErrFoo",
			exp:  false,
		},
		{
			name: "string error",
			err:  fmt.Errorf("normal error"),
			code: "ErrFoo",
			exp:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsErrorCode(tc.err, tc.code)
			if tc.exp != actual {
				t.Errorf("Expected %t but got %t", tc.exp, actual)
			}
		})
	}
}
