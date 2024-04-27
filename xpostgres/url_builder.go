package xpostgres

import (
	"fmt"
	"net/url"
	"os/user"
)

type URLBuilder struct {
	useUnixSocket bool
	host          string
	port          int
	user          string
	password      string
	db            string
	schema        string
	secure        bool
}

func NewURLBuilder() *URLBuilder {
	return &URLBuilder{
		host: "localhost",
		port: 5432,
	}
}

func (b *URLBuilder) UseUnixSocket(use bool) *URLBuilder {
	b.useUnixSocket = use
	return b
}

func (b *URLBuilder) Host(h string) *URLBuilder {
	b.host = h
	return b
}

func (b *URLBuilder) Port(port int) *URLBuilder {
	b.port = port
	return b
}

func (b *URLBuilder) User(u string) *URLBuilder {
	b.user = u
	return b
}

func (b *URLBuilder) Password(password string) *URLBuilder {
	b.password = password
	return b
}

func (b *URLBuilder) DB(db string) *URLBuilder {
	b.db = db
	return b
}

func (b *URLBuilder) Schema(schema string) *URLBuilder {
	b.schema = schema
	return b
}

func (b *URLBuilder) Secure(secure bool) *URLBuilder {
	b.secure = secure
	return b
}

func (b *URLBuilder) Build() string {
	if b.useUnixSocket {
		u, err := user.Current()
		if err != nil {
			return ""
		}
		if b.schema == "" {
			return fmt.Sprintf("postgres:///%s?host=/var/run/postgresql/", u.Username)
		} else {
			return fmt.Sprintf("postgres:///%s?host=/var/run/postgresql/&search_path=%s", u.Username, b.schema)
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
