package types

import (
	"fmt"
	"net/mail"
	xerr "server/xerrors"
	"strings"
)

func NewEmail(email string) (Email, error) {
	var e Email
	if err := e.isValid(email); err != nil {
		return e, fmt.Errorf("invalid email: %w", err)
	}
	e = Email(email)
	return e, nil
}

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) isValid(email string) error {
	// ParseAddress fails to catch 'missing@domain'
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("failed to parse email: %w", err)
	}

	i := strings.Index(email, "@")
	domain := email[i+1:]
	if !strings.Contains(domain, ".") {
		return xerr.EmailDomainMissing
	}

	return nil
}
