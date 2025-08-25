package api

import (
	"encoding/json"
	"net/http"
	mw "server/handler/middleware"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	"server/util"
)

func Login(a i.AuthService, u i.UserService) http.Handler {
	h := login{
		authS: a,
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithNoAuth(), mw.WithMethod(mw.POST))
}

type login struct {
	authS i.AuthService
	userS i.UserService
}

func (l login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	email, err := cred.NewEmail(req.Email)
	if err != nil {
		util.Log.Errorf("invalid email: %v", err)
		SendErrorResponse(w, "Invalid Email format", false)
		return
	}

	emails := []cred.Email{email}

	// Need to add some type of rollback in the case of an error?
	users, err := l.userS.FindUsers(emails)
	if err != nil {
		util.Log.Errorf("failed to find user %v", err)
		SendErrorResponse(w, "Email not found", false)
		return
	}

	if len(users) == 0 {
		util.Log.Errorf("failed to find user %v", err)
		SendErrorResponse(w, "Email not found", false)
		return
	}

	user := users[0]

	pwdBytes := []byte(req.Password)
	if err := user.PwdHash.Compare(pwdBytes); err != nil {
		util.Log.Errorf("invalid password: %v", err)
		SendErrorResponse(w, "Invalid password", false)
		return
	}

	session, err := l.authS.NewSession(user.Id)
	if err != nil {
		util.Log.Errorf("failed to create new user %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	LoginSuccessful(w)
}

func LoginSuccessful(w http.ResponseWriter) {
	SendErrorResponse(w, "", true)
}

func SendErrorResponse(w http.ResponseWriter, message string, noError bool) {
	errorResponse := errorResponse{
		NoError:      noError,
		ErrorMessage: message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse)
}

type loginRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type errorResponse struct {
	NoError      bool   `json:"NoError"`
	ErrorMessage string `json:"ErrorMessage"`
}
