package xpostgres

import (
	"fmt"
	"net/url"
	"os/user"
	"strings"

	"go.olapie.com/x/xurl"
)

type ConnStringBuilder struct {
	useUnixDomainSocket bool
	host                string
	port                int
	user                string
	password            string
	db                  string
	schema              string
	sslMode             string

	query url.Values
}

func NewConnStringBuilder() *ConnStringBuilder {
	return &ConnStringBuilder{
		host:  "localhost",
		port:  5432,
		query: map[string][]string{},
	}
}

func (b *ConnStringBuilder) UseUnixDomainSocket(use bool) *ConnStringBuilder {
	b.useUnixDomainSocket = use
	return b
}

func (b *ConnStringBuilder) Host(h string) *ConnStringBuilder {
	b.host = h
	return b
}

func (b *ConnStringBuilder) Port(port int) *ConnStringBuilder {
	b.port = port
	return b
}

func (b *ConnStringBuilder) User(u string) *ConnStringBuilder {
	b.user = u
	return b
}

func (b *ConnStringBuilder) Password(password string) *ConnStringBuilder {
	b.password = password
	return b
}

func (b *ConnStringBuilder) DB(db string) *ConnStringBuilder {
	b.db = db
	return b
}

func (b *ConnStringBuilder) Schema(schema string) *ConnStringBuilder {
	b.schema = schema
	return b
}

// SSLMode set sslmode, e.g. disabled, verify-ca
func (b *ConnStringBuilder) SSLMode(mode string) *ConnStringBuilder {
	b.sslMode = mode
	return b
}

func (b *ConnStringBuilder) WithQuery(key, value string) *ConnStringBuilder {
	b.query.Set(key, value)
	return b
}

func (b *ConnStringBuilder) Build() string {
	if b.useUnixDomainSocket {
		params := map[string]string{}
		username := b.user
		if username == "" {
			u, err := user.Current()
			if err != nil {
				return ""
			}
			username = u.Username
		}
		params["user"] = username
		port := 5432
		if b.port != 0 {
			port = b.port
		}
		params["port"] = fmt.Sprint(port)
		if b.schema != "" {
			params["search_path"] = b.schema
		}

		if b.db != "" {
			params["database"] = b.db
		}

		s := strings.Builder{}
		for k, v := range params {
			if s.Len() > 0 {
				s.WriteByte(' ')
			}
			s.WriteString(fmt.Sprintf("%s=%v", k, v))
		}
		return s.String()
	}
	host := b.host
	port := b.port
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 5432
	}

	connStr := fmt.Sprintf("%s:%d", host, port)
	if b.db != "" {
		connStr += "/" + b.db
	}
	if b.user == "" {
		connStr = "postgres://" + connStr
	} else {
		if b.password == "" {
			connStr = "postgres://" + b.user + "@" + connStr
		} else {
			connStr = "postgres://" + b.user + ":" + b.password + "@" + connStr
		}
	}
	if b.sslMode != "" {
		b.query.Set("sslmode", b.sslMode)
	}
	if b.schema != "" {
		b.query.Set("search_path", b.schema)
	}
	if len(b.query) == 0 {
		return connStr
	}
	return connStr + "?" + b.query.Encode()
}

func SetParameterInConnString(s string, name, val string) string {
	if strings.Contains(s, "://") {
		updated, err := xurl.SetQuery(s, name, val)
		if err != nil {
			return s
		}
		return updated
	}

	i := strings.Index(s, name+"=")
	if i < 0 {
		return fmt.Sprintf("%s %s=%s", s, name, val)
	}

	j := strings.Index(s[i+len(name)+1:], " ")
	if j < 0 {
		return fmt.Sprintf("%s%s=%s", s[:i], name, val)
	}

	return fmt.Sprintf("%s%s=%s%s", s[:i], name, val, s[j:])
}
