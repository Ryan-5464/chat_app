package handler

import (
	"net/http"
	l "server/logging"
	ss "server/services/authService/session"
)

func checkAuthenticationStatus(r *http.Request) (ss.Session, bool) {
	l.Lgr.DLog("Checking authentication status...")

	session := r.Context().Value("session").(ss.Session)

	emptySession := ss.Session{}
	if session == emptySession {
		l.Lgr.DLog("User not authenticated.")
		return emptySession, false
	}

	l.Lgr.DLog("User is authenticated.")
	return session, true
}
