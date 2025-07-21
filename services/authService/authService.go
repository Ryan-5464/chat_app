package authservice

import (
	"fmt"
	"log"
	skey "server/services/authService/secretKeys"
	sess "server/services/authService/session"
	typ "server/types"
)

func NewAuthService() *AuthService {
	return &AuthService{
		sks: *skey.NewSecretKeyService(300),
	}

}

type AuthService struct {
	sks skey.SecretKeyService
}

func (a *AuthService) ValidateAndRefreshSession(token string) (sess.Session, error) {
	log.Println("dummy validation")
	userId := typ.UserId(1)
	s, err := sess.NewSession(userId, a.sks.CurrentKey())
	if err != nil {
		return sess.Session{}, fmt.Errorf("failed to create new session: %w", err)
	}
	return s, nil
}
