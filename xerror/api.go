package xerror

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code,omitempty"`
	SubCode int    `json:"sub_code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *APIError) String() string {
	return e.Error()
}

func (e *APIError) Error() string {
	if e.Message == "" {
		e.Message = http.StatusText(e.Code)
		if e.Message == "" {
			e.Message = fmt.Sprint(e.Code)
		} else if e.SubCode > 0 {
			e.Message = fmt.Sprintf("%s (%d)", e.Message, e.SubCode)
		}
	}
	return e.Message
}

func (e *APIError) Is(target error) bool {
	if e == target {
		return true
	}

	var t *APIError
	if errors.As(target, &t) {
		return t.Code == e.Code && t.SubCode == e.SubCode && t.Message == e.Message
	}
	return false
}

func NewAPIError(code int, msg string) *APIError {
	return &APIError{
		Code:    code,
		Message: msg,
	}
}

func NewAPIErrorf(code int, format string, a ...any) *APIError {
	msg := fmt.Sprintf(format, a...)
	if msg == "" {
		msg = http.StatusText(code)
	}
	return &APIError{
		Code:    code,
		Message: msg,
	}
}

func NewSubAPIError(code, subCode int, message string) *APIError {
	if code <= 0 {
		panic("invalid code")
	}

	if subCode <= 0 {
		panic("invalid subCode")
	}

	if message == "" {
		message = http.StatusText(code)
	}
	return &APIError{
		Code:    code,
		SubCode: subCode,
		Message: message,
	}
}
