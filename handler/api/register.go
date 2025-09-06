package api

import (
	"encoding/json"
	"log"
	"net/http"
	mw "server/handler/middleware"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	ss "server/services/auth/session"
	"server/util"
)

func Register(a i.AuthService, u i.UserService) http.Handler {
	h := register{
		authS: a,
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithNoAuth(), mw.WithMethod(mw.POST))

}

type register struct {
	authS i.AuthService
	userS i.UserService
}

func (h register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	email, err := cred.NewEmail(req.Email)
	if err != nil {
		util.Log.Errorf("invalid email: %v", err)
		SendErrorResponse(w, "Invalid email format", false)
		return
	}

	pwdBytes := []byte(req.Password)
	pwdHash, err := cred.NewPwdHash(pwdBytes)
	if err != nil {
		util.Log.Errorf("invalid password: %v", err)
		SendErrorResponse(w, "Invalid password format", false)
		return
	}

	log.Println("USername: ", req.Name)
	isValidName := cred.NewUsername(req.Name)
	if !isValidName {
		util.Log.Errorf("invalid username: %v", req.Name)
		SendErrorResponse(w, "Invalid username format", false)
		return
	}

	re := rrequest{
		Name:    req.Name,
		Email:   email,
		PwdHash: pwdHash,
	}

	session, err := h.handleRequest(re)
	if err != nil {
		util.Log.Errorf("failed to register user: %v", err)
		SendErrorResponse(w, "Username already exists!", false)
		return

	}

	http.SetCookie(w, session.Cookie())
	RegisterSuccessful(w)
}

func (h register) handleRequest(req rrequest) (ss.Session, error) {
	util.Log.FunctionInfo()

	user, err := h.userS.NewUser(req.Name, req.Email, req.PwdHash)
	if err != nil {
		return ss.Session{}, err
	}

	util.Log.Info("successfully created new user")

	return h.authS.NewSession(user.Id)
}

func RegisterSuccessful(w http.ResponseWriter) {
	SendErrorResponse(w, "", true)
}

type registerRequest struct {
	Name     string
	Email    string
	Password string
}

type rrequest struct {
	Name    string
	Email   cred.Email
	PwdHash cred.PwdHash
}
