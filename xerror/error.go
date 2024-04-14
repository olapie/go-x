package xerror

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"go.olapie.com/x/xconv"
)

type String string

func (s String) Error() string {
	return string(s)
}

const (
	DBNoRecord String = "database: no records"
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

func New(code int, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if msg == "" {
		msg = http.StatusText(code)
	}
	return &APIError{
		Code:    code,
		Message: msg,
	}
}

func NewSub(code, subCode int, message string) error {
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

func BadRequest(format string, a ...any) error {
	return New(http.StatusBadRequest, format, a...)
}

func Unauthorized(format string, a ...any) error {
	return New(http.StatusUnauthorized, format, a...)
}

func PaymentRequired(format string, a ...any) error {
	return New(http.StatusPaymentRequired, format, a...)
}

func Forbidden(format string, a ...any) error {
	return New(http.StatusForbidden, format, a...)
}

func NotFound(format string, a ...any) error {
	return New(http.StatusNotFound, format, a...)
}

func MethodNotAllowed(format string, a ...any) error {
	return New(http.StatusMethodNotAllowed, format, a...)
}

func NotAcceptable(format string, a ...any) error {
	return New(http.StatusNotAcceptable, format, a...)
}

func ProxyAuthRequired(format string, a ...any) error {
	return New(http.StatusProxyAuthRequired, format, a...)
}

func RequestTimeout(format string, a ...any) error {
	return New(http.StatusRequestTimeout, format, a...)
}

func Conflict(format string, a ...any) error {
	return New(http.StatusConflict, format, a...)
}

func LengthRequired(format string, a ...any) error {
	return New(http.StatusLengthRequired, format, a...)
}

func PreconditionFailed(format string, a ...any) error {
	return New(http.StatusPreconditionFailed, format, a...)
}

func RequestEntityTooLarge(format string, a ...any) error {
	return New(http.StatusRequestEntityTooLarge, format, a...)
}

func RequestURITooLong(format string, a ...any) error {
	return New(http.StatusRequestURITooLong, format, a...)
}

func ExpectationFailed(format string, a ...any) error {
	return New(http.StatusExpectationFailed, format, a...)
}

func Teapot(format string, a ...any) error {
	return New(http.StatusTeapot, format, a...)
}

func MisdirectedRequest(format string, a ...any) error {
	return New(http.StatusMisdirectedRequest, format, a...)
}

func UnprocessableEntity(format string, a ...any) error {
	return New(http.StatusUnprocessableEntity, format, a...)
}

func Locked(format string, a ...any) error {
	return New(http.StatusLocked, format, a...)
}

func TooEarly(format string, a ...any) error {
	return New(http.StatusTooEarly, format, a...)
}

func UpgradeRequired(format string, a ...any) error {
	return New(http.StatusUpgradeRequired, format, a...)
}

func PreconditionRequired(format string, a ...any) error {
	return New(http.StatusPreconditionRequired, format, a...)
}

func TooManyRequests(format string, a ...any) error {
	return New(http.StatusTooManyRequests, format, a...)
}

func RequestHeaderFieldsTooLarge(format string, a ...any) error {
	return New(http.StatusRequestHeaderFieldsTooLarge, format, a...)
}

func UnavailableForLegalReasons(format string, a ...any) error {
	return New(http.StatusUnavailableForLegalReasons, format, a...)
}

func InternalServerError(format string, a ...any) error {
	return New(http.StatusInternalServerError, format, a...)
}

func NotImplemented(format string, a ...any) error {
	return New(http.StatusNotImplemented, format, a...)
}

func BadGateway(format string, a ...any) error {
	return New(http.StatusBadGateway, format, a...)
}

func ServiceUnavailable(format string, a ...any) error {
	return New(http.StatusServiceUnavailable, format, a...)
}

func GatewayTimeout(format string, a ...any) error {
	return New(http.StatusGatewayTimeout, format, a...)
}

func HTTPVersionNotSupported(format string, a ...any) error {
	return New(http.StatusHTTPVersionNotSupported, format, a...)
}

func VariantAlsoNegotiates(format string, a ...any) error {
	return New(http.StatusVariantAlsoNegotiates, format, a...)
}

func InsufficientStorage(format string, a ...any) error {
	return New(http.StatusInsufficientStorage, format, a...)
}

func LoopDetected(format string, a ...any) error {
	return New(http.StatusLoopDetected, format, a...)
}

func NotExtended(format string, a ...any) error {
	return New(http.StatusNotExtended, format, a...)
}

func NetworkAuthenticationRequired(format string, a ...any) error {
	return New(http.StatusNetworkAuthenticationRequired, format, a...)
}

func Wrap(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}
	a = append(a, err)
	return fmt.Errorf(format+":%w", a...)
}

// Cause returns the root cause error
func Cause(err error) error {
	for {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return err
}

func GetCode(err error) int {
	err = Cause(err)
	if err == nil {
		return 0
	}

	if s, ok := err.(interface{ Code() int }); ok {
		return s.Code()
	}

	if s, ok := err.(interface{ GetCode() int }); ok {
		return s.GetCode()
	}

	// ------------
	// int32

	if s, ok := err.(interface{ Code() int32 }); ok {
		return int(s.Code())
	}

	if s, ok := err.(interface{ GetCode() int32 }); ok {
		return int(s.GetCode())
	}

	if s, ok := err.(interface{ StatusCode() int }); ok {
		return s.StatusCode()
	}

	if s, ok := err.(interface{ GetStatusCode() int }); ok {
		return s.GetStatusCode()
	}

	if s, ok := err.(interface{ Status() int }); ok {
		return s.Status()
	}

	if s, ok := err.(interface{ GetStatus() int }); ok {
		return s.GetStatus()
	}

	// ------------
	// int32

	if s, ok := err.(interface{ StatusCode() int32 }); ok {
		return int(s.StatusCode())
	}

	if s, ok := err.(interface{ GetStatusCode() int32 }); ok {
		return int(s.GetStatusCode())
	}

	if s, ok := err.(interface{ Status() int32 }); ok {
		return int(s.Status())
	}

	if s, ok := err.(interface{ GetStatus() int32 }); ok {
		return int(s.GetStatus())
	}

	v := reflect.ValueOf(xconv.Indirect(err))
	t := v.Type()
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			switch ft.Name {
			case "Code", "Status", "StatusCode", "ErrorCode":
				fv := v.Field(i)
				if fv.CanInt() {
					return int(fv.Int())
				}

				if fv.CanUint() {
					return int(fv.Uint())
				}

				return 0
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Kind() != reflect.String {
				continue
			}
			switch k.String() {
			case "Code", "Status", "StatusCode", "ErrorCode":
				vv := v.MapIndex(k)
				if vv.CanInt() {
					return int(vv.Int())
				}

				if vv.CanUint() {
					return int(vv.Uint())
				}

				return 0
			}
		}
	}
	return 0
}
