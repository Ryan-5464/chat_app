package middleware

import (
	"net/http"
	"server/util"
)

func NewNoAuthMW() *noAuthMW {
	return &noAuthMW{}
}

type noAuthMW struct{}

func (a *noAuthMW) Bind(next http.Handler) http.Handler {
	util.Log.FunctionInfo()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Log.FunctionInfo()

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err != http.ErrNoCookie {
				util.Log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if cookie != nil {
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	})
}
