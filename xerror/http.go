package xerror

import "net/http"

func BadRequest(format string, a ...any) error {
	return NewAPIError(http.StatusBadRequest, format, a...)
}

func Unauthorized(format string, a ...any) error {
	return NewAPIError(http.StatusUnauthorized, format, a...)
}

func PaymentRequired(format string, a ...any) error {
	return NewAPIError(http.StatusPaymentRequired, format, a...)
}

func Forbidden(format string, a ...any) error {
	return NewAPIError(http.StatusForbidden, format, a...)
}

func NotFound(format string, a ...any) error {
	return NewAPIError(http.StatusNotFound, format, a...)
}

func MethodNotAllowed(format string, a ...any) error {
	return NewAPIError(http.StatusMethodNotAllowed, format, a...)
}

func NotAcceptable(format string, a ...any) error {
	return NewAPIError(http.StatusNotAcceptable, format, a...)
}

func ProxyAuthRequired(format string, a ...any) error {
	return NewAPIError(http.StatusProxyAuthRequired, format, a...)
}

func RequestTimeout(format string, a ...any) error {
	return NewAPIError(http.StatusRequestTimeout, format, a...)
}

func Conflict(format string, a ...any) error {
	return NewAPIError(http.StatusConflict, format, a...)
}

func LengthRequired(format string, a ...any) error {
	return NewAPIError(http.StatusLengthRequired, format, a...)
}

func PreconditionFailed(format string, a ...any) error {
	return NewAPIError(http.StatusPreconditionFailed, format, a...)
}

func RequestEntityTooLarge(format string, a ...any) error {
	return NewAPIError(http.StatusRequestEntityTooLarge, format, a...)
}

func RequestURITooLong(format string, a ...any) error {
	return NewAPIError(http.StatusRequestURITooLong, format, a...)
}

func ExpectationFailed(format string, a ...any) error {
	return NewAPIError(http.StatusExpectationFailed, format, a...)
}

func Teapot(format string, a ...any) error {
	return NewAPIError(http.StatusTeapot, format, a...)
}

func MisdirectedRequest(format string, a ...any) error {
	return NewAPIError(http.StatusMisdirectedRequest, format, a...)
}

func UnprocessableEntity(format string, a ...any) error {
	return NewAPIError(http.StatusUnprocessableEntity, format, a...)
}

func Locked(format string, a ...any) error {
	return NewAPIError(http.StatusLocked, format, a...)
}

func TooEarly(format string, a ...any) error {
	return NewAPIError(http.StatusTooEarly, format, a...)
}

func UpgradeRequired(format string, a ...any) error {
	return NewAPIError(http.StatusUpgradeRequired, format, a...)
}

func PreconditionRequired(format string, a ...any) error {
	return NewAPIError(http.StatusPreconditionRequired, format, a...)
}

func TooManyRequests(format string, a ...any) error {
	return NewAPIError(http.StatusTooManyRequests, format, a...)
}

func RequestHeaderFieldsTooLarge(format string, a ...any) error {
	return NewAPIError(http.StatusRequestHeaderFieldsTooLarge, format, a...)
}

func UnavailableForLegalReasons(format string, a ...any) error {
	return NewAPIError(http.StatusUnavailableForLegalReasons, format, a...)
}

func InternalServerError(format string, a ...any) error {
	return NewAPIError(http.StatusInternalServerError, format, a...)
}

func NotImplemented(format string, a ...any) error {
	return NewAPIError(http.StatusNotImplemented, format, a...)
}

func BadGateway(format string, a ...any) error {
	return NewAPIError(http.StatusBadGateway, format, a...)
}

func ServiceUnavailable(format string, a ...any) error {
	return NewAPIError(http.StatusServiceUnavailable, format, a...)
}

func GatewayTimeout(format string, a ...any) error {
	return NewAPIError(http.StatusGatewayTimeout, format, a...)
}

func HTTPVersionNotSupported(format string, a ...any) error {
	return NewAPIError(http.StatusHTTPVersionNotSupported, format, a...)
}

func VariantAlsoNegotiates(format string, a ...any) error {
	return NewAPIError(http.StatusVariantAlsoNegotiates, format, a...)
}

func InsufficientStorage(format string, a ...any) error {
	return NewAPIError(http.StatusInsufficientStorage, format, a...)
}

func LoopDetected(format string, a ...any) error {
	return NewAPIError(http.StatusLoopDetected, format, a...)
}

func NotExtended(format string, a ...any) error {
	return NewAPIError(http.StatusNotExtended, format, a...)
}

func NetworkAuthenticationRequired(format string, a ...any) error {
	return NewAPIError(http.StatusNetworkAuthenticationRequired, format, a...)
}
