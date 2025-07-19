package errors

import (
	e "errors"
)

var (
	InvalidEmail      error = e.New("invalid email address")
	JWEGenerationFail error = e.New("failed to generate JWE")
	NoUpperCaseChar   error = e.New("no uppercase found")
	NoLowerCaseChar   error = e.New("no lowercase found")
	NoDigitChar       error = e.New("no digit found")
	NoSymbolChar      error = e.New("no symbol found")
	InvalidPwd        error = e.New("password and hash do not match")
	PwdHashFail       error = e.New("failed to hash password")
	PwdTooShort       error = e.New("password too short")
)
