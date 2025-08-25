package api

import (
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

	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	emailStr := r.FormValue("Email")
	passwordStr := r.FormValue("Password")

	email, err := cred.NewEmail(emailStr)
	if err != nil {
		util.Log.Errorf("invalid email: %v", err)
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	pwdBytes := []byte(passwordStr)
	pwdHash, err := cred.NewPwdHash(pwdBytes)
	if err != nil {
		util.Log.Errorf("invalid password: %v", err)
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	req := rrequest{
		Name:    name,
		Email:   email,
		PwdHash: pwdHash,
	}

	session, err := h.handleRequest(req)
	if err != nil {
		util.Log.Errorf("failed to register user: %v", err)
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return

	}

	http.SetCookie(w, session.Cookie())

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
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

type rrequest struct {
	Name    string
	Email   cred.Email
	PwdHash cred.PwdHash
}
