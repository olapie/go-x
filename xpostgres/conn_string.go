package xpostgres

import (
	"fmt"
	"go.olapie.com/x/xurl"
	"net/url"
	"os/user"
	"strings"
)

type ConnStringBuilder struct {
	useUnixDomainSocket bool
	host                string
	port                int
	user                string
	password            string
	db                  string
	schema              string
	secure              bool
}

func NewConnStringBuilder() *ConnStringBuilder {
	return &ConnStringBuilder{
		host: "localhost",
		port: 5432,
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

func (b *ConnStringBuilder) Secure(secure bool) *ConnStringBuilder {
	b.secure = secure
	return b
}

func (b *ConnStringBuilder) Build() string {
	if b.useUnixDomainSocket {
		username := b.user
		if username == "" {
			u, err := user.Current()
			if err != nil {
				return ""
			}
			username = u.Username
		}
		port := 5432
		if b.port != 0 {
			port = b.port
		}
		if b.schema == "" {
			return fmt.Sprintf("host=/tmp port=%d user=%s", port, username)
		} else {
			return fmt.Sprintf("host=/tmp port=%d user=%s search_path=%s", port, username, b.schema)
		}
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
	query := url.Values{}
	if !b.secure {
		query.Add("sslmode", "disable")
	}
	if b.schema != "" {
		query.Add("search_path", b.schema)
	}
	if len(query) == 0 {
		return connStr
	}
	return connStr + "?" + query.Encode()
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
