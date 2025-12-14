package xmobile

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.olapie.com/x/xconv"
	"go.olapie.com/x/xtime"
)

type Uptimer interface {
	Uptime() int64
}

const (
	KeyDeviceID = "mobile_device_id"
)

func GetDeviceID(m SecretManager) string {
	id := m.Get(KeyDeviceID)
	if id == "" {
		id = uuid.NewString()
		m.Set(KeyDeviceID, id)
	}
	return id
}

func SetTimeZone(name string, offset int) {
	time.Local = xtime.ToLocation(name, offset)
}

func GetTimeZoneOffset() int {
	_, o := time.Now().Zone()
	return o
}

func GetTimeZoneName() string {
	n, _ := time.Now().Zone()
	return n
}

func NewUUID() string {
	return uuid.New().String()
}

type Now interface {
	Now() int64
}

func SmartLen(s string) int {
	n := 0
	for _, c := range s {
		if c <= 255 {
			n++
		} else {
			n += 2
		}
	}

	return n
}

var whitespaceRegexp = regexp.MustCompile(`[ \t\n\r]+`)

// SquishString returns the string
// first removing all whitespace on both ends of the string,
// and then changing remaining consecutive whitespace groups into one space each.
func SquishString[T ~string](s T) T {
	str := strings.TrimSpace(string(s))
	str = whitespaceRegexp.ReplaceAllString(str, " ")
	return T(str)
}

type Handler interface {
	SaveSecret(name, value string) bool
	DeleteSecret(name string) bool
	GetSecret(name string) string
	NeedSignIn()
}

type AuthErrorChecker struct {
	h Handler
}

func NewAuthErrorChecker(h Handler) *AuthErrorChecker {
	return &AuthErrorChecker{
		h: h,
	}
}

func (c *AuthErrorChecker) Check(err error) {
	if err == nil {
		return
	}

	code := ToError(err).Code
	if code == http.StatusUnauthorized {
		go c.h.NeedSignIn()
	}
}

func GetSizeString(n int64) string {
	return xconv.SizeToHumanReadable(n)
}

type Context struct {
	ctx context.Context
}

func NewContext() *Context {
	return &Context{
		ctx: context.Background(),
	}
}

func (c *Context) Get() context.Context {
	return c.ctx
}

func (c *Context) Set(ctx context.Context) {
	c.ctx = ctx
}
