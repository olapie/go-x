package xerror

import "fmt"

func PanicWithMessages(msgAndArgs ...any) {
	n := len(msgAndArgs)
	switch n {
	case 0:
		panic("")
	case 1:
		panic(msgAndArgs[0])
	default:
		if format, ok := msgAndArgs[0].(string); ok {
			panic(fmt.Sprintf(format, msgAndArgs[1:]...))
		}
		panic(fmt.Sprint(msgAndArgs...))
	}
}

// Recover recovers from panic and assign message to outErr
// outErr usually is a pointer to return error
// E.g.
//
//	func doSomething() (err error) {
//	    defer Recover(&err)
//	    ...
//	}
func Recover(outErr *error) {
	if v := recover(); v != nil {
		if err, ok := v.(error); ok {
			*outErr = err
		} else {
			*outErr = fmt.Errorf("panic: %v", v)
		}
	}
}
