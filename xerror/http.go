package xerror

import "net/http"

func BadRequest(format string, a ...any) error {
	return NewAPIErrorf(http.StatusBadRequest, format, a...)
}

func Unauthorized(format string, a ...any) error {
	return NewAPIErrorf(http.StatusUnauthorized, format, a...)
}

func PaymentRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusPaymentRequired, format, a...)
}

func Forbidden(format string, a ...any) error {
	return NewAPIErrorf(http.StatusForbidden, format, a...)
}

func NotFound(format string, a ...any) error {
	return NewAPIErrorf(http.StatusNotFound, format, a...)
}

func MethodNotAllowed(format string, a ...any) error {
	return NewAPIErrorf(http.StatusMethodNotAllowed, format, a...)
}

func NotAcceptable(format string, a ...any) error {
	return NewAPIErrorf(http.StatusNotAcceptable, format, a...)
}

func ProxyAuthRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusProxyAuthRequired, format, a...)
}

func RequestTimeout(format string, a ...any) error {
	return NewAPIErrorf(http.StatusRequestTimeout, format, a...)
}

func Conflict(format string, a ...any) error {
	return NewAPIErrorf(http.StatusConflict, format, a...)
}

func LengthRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusLengthRequired, format, a...)
}

func PreconditionFailed(format string, a ...any) error {
	return NewAPIErrorf(http.StatusPreconditionFailed, format, a...)
}

func RequestEntityTooLarge(format string, a ...any) error {
	return NewAPIErrorf(http.StatusRequestEntityTooLarge, format, a...)
}

func RequestURITooLong(format string, a ...any) error {
	return NewAPIErrorf(http.StatusRequestURITooLong, format, a...)
}

func ExpectationFailed(format string, a ...any) error {
	return NewAPIErrorf(http.StatusExpectationFailed, format, a...)
}

func Teapot(format string, a ...any) error {
	return NewAPIErrorf(http.StatusTeapot, format, a...)
}

func MisdirectedRequest(format string, a ...any) error {
	return NewAPIErrorf(http.StatusMisdirectedRequest, format, a...)
}

func UnprocessableEntity(format string, a ...any) error {
	return NewAPIErrorf(http.StatusUnprocessableEntity, format, a...)
}

func Locked(format string, a ...any) error {
	return NewAPIErrorf(http.StatusLocked, format, a...)
}

func TooEarly(format string, a ...any) error {
	return NewAPIErrorf(http.StatusTooEarly, format, a...)
}

func UpgradeRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusUpgradeRequired, format, a...)
}

func PreconditionRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusPreconditionRequired, format, a...)
}

func TooManyRequests(format string, a ...any) error {
	return NewAPIErrorf(http.StatusTooManyRequests, format, a...)
}

func RequestHeaderFieldsTooLarge(format string, a ...any) error {
	return NewAPIErrorf(http.StatusRequestHeaderFieldsTooLarge, format, a...)
}

func UnavailableForLegalReasons(format string, a ...any) error {
	return NewAPIErrorf(http.StatusUnavailableForLegalReasons, format, a...)
}

func InternalServerError(format string, a ...any) error {
	return NewAPIErrorf(http.StatusInternalServerError, format, a...)
}

func NotImplemented(format string, a ...any) error {
	return NewAPIErrorf(http.StatusNotImplemented, format, a...)
}

func BadGateway(format string, a ...any) error {
	return NewAPIErrorf(http.StatusBadGateway, format, a...)
}

func ServiceUnavailable(format string, a ...any) error {
	return NewAPIErrorf(http.StatusServiceUnavailable, format, a...)
}

func GatewayTimeout(format string, a ...any) error {
	return NewAPIErrorf(http.StatusGatewayTimeout, format, a...)
}

func HTTPVersionNotSupported(format string, a ...any) error {
	return NewAPIErrorf(http.StatusHTTPVersionNotSupported, format, a...)
}

func VariantAlsoNegotiates(format string, a ...any) error {
	return NewAPIErrorf(http.StatusVariantAlsoNegotiates, format, a...)
}

func InsufficientStorage(format string, a ...any) error {
	return NewAPIErrorf(http.StatusInsufficientStorage, format, a...)
}

func LoopDetected(format string, a ...any) error {
	return NewAPIErrorf(http.StatusLoopDetected, format, a...)
}

func NotExtended(format string, a ...any) error {
	return NewAPIErrorf(http.StatusNotExtended, format, a...)
}

func NetworkAuthenticationRequired(format string, a ...any) error {
	return NewAPIErrorf(http.StatusNetworkAuthenticationRequired, format, a...)
}
