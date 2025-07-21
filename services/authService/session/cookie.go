package session

import (
	"fmt"
	"net/http"
	tkn "server/services/authService/jwetoken"
	skey "server/services/authService/secretKeys"
	typ "server/types"
)

func NewSession(userId typ.UserId, key skey.SecretKey) (Session, error) {
	jwe, err := tkn.NewJWE(userId, key)
	if err != nil {
		return Session{}, fmt.Errorf("jwe generation failed: %w", err)
	}

	s := Session{}
	s.SetCookie(NewCookie(jwe))
	s.SetUserId(userId)
	return s, nil
}

func NewCookie(jwe tkn.JWE) http.Cookie {
	return http.Cookie{
		Name:     "session_token",
		Value:    jwe.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  jwe.TokenExpiry(),
	}
}

type Session struct {
	cookie http.Cookie
	userId typ.UserId
}

func (s Session) SetCookie(c http.Cookie) {
	s.cookie = c
}

func (s Session) SetUserId(u typ.UserId) {
	s.userId = u
}

func (s Session) Cookie() http.Cookie {
	return s.cookie
}

func (s Session) UserId() typ.UserId {
	return s.userId
}
