package errors

import (
	e "errors"
)

var (
	EmailDomainMissing error = e.New("domain missing after @")
	InitDbServiceFail  error = e.New("no database service initialized")
	NoUpperCaseChar    error = e.New("no uppercase found")
	NoLowerCaseChar    error = e.New("no lowercase found")
	NoDigitChar        error = e.New("no digit found")
	NoSymbolChar       error = e.New("no symbol found")
	PwdTooShort        error = e.New("password too short")
	UserIdNotFound     error = e.New("userid missing")
)
