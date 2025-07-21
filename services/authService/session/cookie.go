package session

import (
	"fmt"
	"net/http"
	tkn "server/services/authService/jwetoken"
	skey "server/services/authService/secretKeys"
	typ "server/types"
	"time"
)

func NewSession(userId typ.UserId, key skey.SecretKey) (Session, error) {
	jwe, err := tkn.NewJWE(userId, key)
	if err != nil {
		return Session{}, fmt.Errorf("jwe generation failed: %w", err)
	}

	c := http.Cookie{
		Name:     "session_token",
		Value:    jwe.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  jwe.TokenExpiry(),
	}

	s := Session{
		cookie: c,
		userId: userId,
	}

	return s, nil
}

type Session struct {
	cookie http.Cookie
	userId typ.UserId
}

func (s Session) Cookie() http.Cookie {
	return s.cookie
}

func (s Session) UserId() typ.UserId {
	return s.userId
}

func (s Session) TokenExpiry() time.Time {
	return s.cookie.Expires
}

func (s Session) JWEToken() string {
	return s.cookie.Value
}

func (s Session) Name() string {
	return s.cookie.Name
}
