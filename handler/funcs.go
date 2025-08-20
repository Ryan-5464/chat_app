package handler

import (
	"encoding/json"
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

func SendJSONResponse(w http.ResponseWriter, responseDTO any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		// h.lgr.LogError(fmt.Errorf("failed to encode JSON response, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
}
