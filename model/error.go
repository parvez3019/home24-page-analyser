package model

import (
	"fmt"
)

// Err represents a model for custom error
type Err struct {
	ErrorCode  ErrorCode `json:"error"`
	Message    string    `json:"message"`
	StatusCode int       `json:"statusCode"`
}

func NewError(code ErrorCode, msg string, statusCode int) *Err {
	return &Err{
		ErrorCode:  code,
		Message:    msg,
		StatusCode: statusCode,
	}
}

//interface error
func (e *Err) Error() string {
	return fmt.Sprintf("Error %s, %s", e.ErrorCode, e.Message)
}

// ErrorCode is used on the top level of errors and is used as a discriminator.
type ErrorCode string

const (
	ErrorCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrorCodeSomethingWentWrong ErrorCode = "SOMETHING_WENT_WRONG"
)
