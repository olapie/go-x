package xmobile

import (
	"errors"
	"fmt"
	"reflect"

	"go.olapie.com/x/xerror"
)

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var _ error = (*Error)(nil)

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %d", e.Message, e.Code)
}

func ToError(err error) *Error {
	if err == nil {
		return nil
	}

	if v := reflect.ValueOf(err); !v.IsValid() || v.IsZero() {
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		return e
	}

	code := xerror.GetCode(err)
	return NewError(code, err.Error())
}
