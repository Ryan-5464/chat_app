package handler

import (
	"errors"
	"log"
	"net/http"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	"text/template"
)

func NewRegisterHandler(l i.Logger, a i.AuthService, u i.UserService) *RegisterHandler {
	return &RegisterHandler{
		lgr:   l,
		authS: a,
		userS: u,
	}
}

type RegisterHandler struct {
	authS i.AuthService
	userS i.UserService
	lgr   i.Logger
}

func (rh *RegisterHandler) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	rh.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cookieFound bool
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			cookieFound = true
		} else {
			cookieFound = false
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if cookieFound {
		token := cookie.Value
		session, err := l.authS.ValidateAndRefreshSession(token)
		if err != nil {
			log.Println("error validating or refreshing session", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, session.Cookie())

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./static/pages/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (rh *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	rh.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("failed to parse form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	emailStr := r.FormValue("Email")
	passwordStr := r.FormValue("Password")

	email, err := cred.NewEmail(emailStr)
	if err != nil {
		log.Printf("invalid email: %v", err)
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	pwdBytes := []byte(passwordStr)
	pwdHash, err := cred.NewPwdHash(pwdBytes)
	if err != nil {
		log.Printf("invalid password: %v", err)
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	usrE := ent.User{
		Email:   email,
		PwdHash: pwdHash,
	}

	// Need to add some type of rollback in the case of an error
	usr, err := rh.userS.NewUser(usrE)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("successfully created new user")

	session, err := rh.authS.NewSession(usr.Id)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
