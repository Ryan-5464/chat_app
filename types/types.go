package types

import (
	"net/mail"
	xerr "server/xerrors"
	"strings"
)

func NewEmail(email string) (Email, error) {
	var e Email
	if !e.valid(email) {
		return e, xerr.InvalidEmail
	}
	e = Email(email)
	return e, nil
}

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) valid(email string) bool {
	// ParseAddress fails to catch 'missing@domain'
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	i := strings.Index(email, "@")
	domain := email[i+1:]
	return strings.Contains(domain, ".")
}

type UserId int64

func (u UserId) Int64() int64 {
	return int64(u)
}
