package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	i "server/interfaces"
	ss "server/services/authService/session"
)

func NewAuthMiddleware(l i.Logger, a i.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		lgr:   l,
		authS: a,
	}
}

type AuthMiddleware struct {
	lgr   i.Logger
	authS i.AuthService
}

func (a *AuthMiddleware) AttachTo(next http.Handler) http.Handler {
	a.lgr.LogFunctionInfo()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.lgr.LogFunctionInfo()
		log.Println("AUTHENTICATING...")

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		var session ss.Session
		if cookie != nil {
			token := cookie.Value
			session, err = a.authS.ValidateAndRefreshSession(token)
			if err != nil {
				log.Println("error validating or refreshing session", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, session.Cookie())
		}

		ctx := context.WithValue(r.Context(), "session", session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
