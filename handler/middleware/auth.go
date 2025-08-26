package middleware

import (
	"context"
	"errors"
	"net/http"
	"server/handler/ctxutil"
	"server/handler/status"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
)

func NewAuthMW(a i.AuthService) *authMW {
	return &authMW{authS: a}
}

type authMW struct {
	authS i.AuthService
}

func (a *authMW) Bind(next http.Handler) http.Handler {
	util.Log.FunctionInfo()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Log.FunctionInfo()

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				util.Log.Infof("%s => redirecting to landing page...", status.NoSessionCookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				util.Log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		token := cookie.Value
		session, err := a.authS.ValidateAndRefreshSession(token)
		if err != nil {
			util.Log.Errorf("%s => error: %v", status.InvalidSession, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, session.Cookie())

		ctx := context.WithValue(r.Context(), ctxutil.SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkAuthenticationStatus(r *http.Request) (ss.Session, bool) {
	util.Log.Dbug("Checking authentication status...")

	session := r.Context().Value("session").(ss.Session)

	emptySession := ss.Session{}
	if session == emptySession {
		util.Log.Dbug("User not authenticated.")
		return emptySession, false
	}

	util.Log.Dbug("User is authenticated.")
	return session, true
}
