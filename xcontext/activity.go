package xcontext

import (
	"fmt"
	"log/slog"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"

	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xsession"
	"go.olapie.com/x/xtype"
)

const (
	ErrNotExist xerror.String = "Activity does not exist"
)

type HeaderTypes interface {
	~map[string][]string | ~map[string]string
}

type Activity struct {
	name   string
	header map[string][]string

	//Session is only available in incoming context, may be nil if session is not enabled
	session *xsession.Session
	userID  xtype.UserIDInterface
}

var (
	typeMapStringToStringSlice = reflect.TypeOf(map[string][]string(nil))
	typeMapStringToString      = reflect.TypeOf(map[string]string(nil))
)

func NewActivity[H HeaderTypes](name string, header H) *Activity {
	a := &Activity{
		name:   name,
		header: make(map[string][]string),
	}

	if header == nil {
		panic("header is nil")
	}

	hv := reflect.ValueOf(header)
	if hv.CanConvert(typeMapStringToStringSlice) {
		a.header = copyHeader(hv.Convert(typeMapStringToStringSlice).Interface().(map[string][]string))
	} else if hv.CanConvert(typeMapStringToString) {
		a.header = copyHeader(hv.Convert(typeMapStringToString).Interface().(map[string]string))
	} else {
		panic(fmt.Sprintf("unsupported header type: %T", header))
	}
	return a
}

func (a *Activity) Name() string {
	return a.name
}

func (a *Activity) Session() *xsession.Session {
	return a.session
}

func (a *Activity) UserID() xtype.UserIDInterface {
	return a.userID
}

func (a *Activity) SetUserID(id xtype.UserIDInterface) {
	a.userID = id
}

func (a *Activity) Set(key string, value string) {
	key = strings.ToLower(key)
	a.header[key] = []string{value}
}

func (a *Activity) Append(key string, value string) {
	key = strings.ToLower(key)
	a.header[key] = append(a.header[key], value)
}

func (a *Activity) Get(key string) string {
	if v := a.get(key); v != "" {
		return v
	}

	if v := a.get(strings.ToLower(key)); v != "" {
		return v
	}

	if v := a.get(textproto.CanonicalMIMEHeaderKey(key)); v != "" {
		return v
	}

	return ""
}

func (a *Activity) get(key string) string {
	v, _ := a.header[strings.ToLower(key)]
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

func (a *Activity) GetAppID() string {
	return a.Get(xhttpheader.KeyAppID)
}

func (a *Activity) SetAppID(id string) {
	a.Set(xhttpheader.KeyAppID, id)
}

func (a *Activity) GetTraceID() string {
	return a.Get(xhttpheader.KeyTraceID)
}

func (a *Activity) SetTraceID(id string) {
	a.Set(xhttpheader.KeyTraceID, id)
}

func (a *Activity) GetClientID() string {
	return a.Get(xhttpheader.KeyClientID)
}

func (a *Activity) SetClientID(id string) {
	a.Set(xhttpheader.KeyClientID, id)
}

func (a *Activity) GetAuthorization() string {
	return a.Get(xhttpheader.KeyAuthorization)
}

func (a *Activity) SetAuthorization(auth string) {
	a.Set(xhttpheader.KeyAuthorization, auth)
}

func (a *Activity) GetRequestTimeout() int {
	s := a.Get(xhttpheader.KeyRequestTimeout)
	if s == "" {
		return 0
	}
	t, err := strconv.Atoi(s)
	if err != nil || t < 0 {
		slog.Error("invalid Request-Timeout: " + s)
		return 0
	}
	return t
}

func (a *Activity) SetRequestTimeout(seconds int) {
	if seconds > 0 {
		a.Set(xhttpheader.KeyRequestTimeout, strconv.Itoa(seconds))
	}
}

func CopyActivityHeader[H HeaderTypes](dest H, a *Activity) {

	if setter, ok := any(dest).(interface {
		Set(key string, values ...string)
	}); ok {
		for k, v := range a.header {
			setter.Set(k, v...)
		}
		return
	}

	if setter, ok := any(dest).(interface {
		Set(key string, values []string)
	}); ok {
		for k, v := range a.header {
			setter.Set(k, v)
		}
		return
	}

	if setter, ok := any(dest).(interface {
		Set(key string, value string)
	}); ok {
		for k, v := range a.header {
			if len(v) > 0 {
				setter.Set(k, v[0])
			}
		}
		return
	}

	if hm, ok := any(dest).(map[string][]string); ok {
		for k, v := range a.header {
			hm[k] = v
		}
		return
	}

	if hs, ok := any(dest).(map[string]string); ok {
		for k, v := range a.header {
			if len(v) != 0 {
				hs[k] = v[0]
			}
		}
		return
	}

	panic(fmt.Sprintf("unsupported header type: %T", dest))
}

func copyHeader[H map[string][]string | map[string]string](h H) map[string][]string {
	res := make(map[string][]string, len(h))
	if m, ok := any(h).(map[string][]string); ok {
		for k, v := range m {
			res[strings.ToLower(k)] = v
		}
		return res
	}

	m := any(h).(map[string]string)
	for k, v := range m {
		res[strings.ToLower(k)] = []string{v}
	}
	return res
}
