package authservice

import (
	"fmt"
	i "server/interfaces"
	skey "server/services/authService/secretKeys"
	sess "server/services/authService/session"
	typ "server/types"
)

func NewAuthService(lgr i.Logger) *AuthService {
	return &AuthService{
		lgr: lgr,
		sks: *skey.NewSecretKeyService(300),
	}

}

type AuthService struct {
	lgr i.Logger
	sks skey.SecretKeyService
}

func (a *AuthService) ValidateAndRefreshSession(token string) (sess.Session, error) {
	a.lgr.LogFunctionInfo()
	userId := typ.UserId(1)
	s, err := a.NewSession(userId)
	if err != nil {
		return sess.Session{}, fmt.Errorf("failed to create new session: %w", err)
	}
	return s, nil
}

func (a *AuthService) NewSession(userId typ.UserId) (sess.Session, error) {
	a.lgr.LogFunctionInfo()
	s, err := sess.NewSession(userId, a.sks.CurrentKey())
	if err != nil {
		return sess.Session{}, fmt.Errorf("failed to create new session: %w", err)
	}
	return s, nil
}
