package credentials

import (
	"fmt"
	xerr "server/xerrors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func NewPwdHash(pwd []byte) (PwdHash, error) {
	var p PwdHash
	if p.alreadyHashed(pwd) {
		return PwdHash(pwd), nil
	}

	if len(pwd) < 12 {
		return nil, xerr.PwdTooShort
	}

	if err := p.valid(pwd); err != nil {
		return nil, fmt.Errorf("password valdiation failed %w", err)
	}

	pwdHash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}
	return PwdHash(pwdHash), nil
}

type PwdHash []byte

func (p PwdHash) String() string {
	return string(p)
}

func (p PwdHash) Compare(pwd []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(p), pwd)
	if err != nil {
		return fmt.Errorf("password hash comparison failed: %w", err)
	}
	return nil
}

func (p PwdHash) alreadyHashed(pwd []byte) bool {
	s := string(pwd)
	if len(s) != 60 {
		return false
	}

	return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
}

func (p PwdHash) valid(pwd []byte) error {
	var hasUpper, hasLower, hasDigit, hasSymbol bool

	for _, ch := range string(pwd) {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return xerr.NoUpperCaseChar
	}
	if !hasLower {
		return xerr.NoLowerCaseChar
	}
	if !hasDigit {
		return xerr.NoDigitChar
	}
	if !hasSymbol {
		return xerr.NoSymbolChar
	}
	return nil
}
