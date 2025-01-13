package xtype

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	nickRegexp     = regexp.MustCompile("^[^ \n\r\t\f][^\n\r\t\f]{0,28}[^ \n\r\t\f]$")
	usernameRegexp = regexp.MustCompile("^[a-zA-Z][\\w\\.]{1,19}$")
)

type Account interface {
	AccountType() string
	String() string
}

type Username string

func (u Username) AccountType() string {
	return "username"
}

func (u Username) IsValid() bool {
	return usernameRegexp.MatchString(string(u))
}

func (u Username) Normalize() Username {
	s := string(u)
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return Username(s)
}

func (u Username) String() string {
	return (string)(u)
}

type EmailAddress string

func (e EmailAddress) AccountType() string {
	return "email_address"
}

func (e EmailAddress) Normalize() EmailAddress {
	s := string(e)
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return EmailAddress(s)
}

func (e EmailAddress) IsValid() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

func (e EmailAddress) String() string {
	return (string)(e)
}

type Nickname string

func (n Nickname) IsValid() bool {
	return nickRegexp.MatchString(string(n))
}

func (n Nickname) Normalize() Nickname {
	s := string(n)
	s = strings.TrimSpace(s)
	return Nickname(s)
}

func RandomEmailAddress() EmailAddress {
	s := RandomID().Base62() + "@" + RandomID().Base62() + ".com"
	return EmailAddress(s)
}

func RandomNickname() Nickname {
	return Nickname(RandomID().Base62())
}
