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
	name          string
	headersMulti  map[string][]string
	headersSingle map[string]string

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
		name: name,
	}

	if header == nil {
		panic("header is nil")
	}

	if hm, ok := any(header).(map[string][]string); ok {
		a.headersMulti = hm
	} else if hs, ok := any(header).(map[string]string); ok {
		a.headersSingle = hs
	} else {
		hv := reflect.ValueOf(header)
		if hv.CanConvert(typeMapStringToStringSlice) {
			a.headersMulti = hv.Convert(typeMapStringToStringSlice).Interface().(map[string][]string)
		} else if hv.CanConvert(typeMapStringToString) {
			a.headersSingle = hv.Convert(typeMapStringToString).Interface().(map[string]string)
		} else {
			panic(fmt.Sprintf("unsupported header type: %T", header))
		}
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
	if a.headersMulti != nil {
		a.headersMulti[key] = append(a.headersMulti[key], value)
	} else {
		a.headersSingle[key] = value
	}
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
	if a.headersMulti != nil {
		l := a.headersMulti[key]
		if len(l) != 0 {
			return l[0]
		}
		return ""
	}

	if v, ok := a.headersSingle[key]; ok {
		return v
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
	if hm, ok := any(dest).(map[string][]string); ok {
		if a.headersMulti != nil {
			for k, v := range a.headersMulti {
				hm[k] = v
			}
		} else {
			for k, v := range a.headersSingle {
				hm[k] = append(hm[k], v)
			}
		}
		return
	}

	if hs, ok := any(dest).(map[string]string); ok {
		if a.headersMulti != nil {
			for k, v := range a.headersMulti {
				if len(v) != 0 {
					hs[k] = v[0]
				}
			}
		} else {
			for k, v := range a.headersSingle {
				hs[k] = v
			}
		}
		return
	}

	panic(fmt.Sprintf("unsupported header type: %T", dest))
}
