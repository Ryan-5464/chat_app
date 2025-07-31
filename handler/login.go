package handler

import (
	"log"
	"net/http"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	ss "server/services/authService/session"
	"text/template"
)

func NewLoginHandler(l i.Logger, a i.AuthService, u i.UserService) *LoginHandler {
	return &LoginHandler{
		lgr:   l,
		authS: a,
		userS: u,
	}
}

type LoginHandler struct {
	authS i.AuthService
	userS i.UserService
	lgr   i.Logger
}

func (l *LoginHandler) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	l.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(ss.Session)
	emptySession := ss.Session{}
	if session != emptySession {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./static/pages/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (l *LoginHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	l.lgr.LogFunctionInfo()

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

	usrE := ent.User{
		Email: email,
	}

	// Need to add some type of rollback in the case of an error?
	usr, err := l.userS.FindUser(usrE)
	if err != nil {
		log.Printf("failed to find user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := usr.PwdHash.Compare(pwdBytes); err != nil {
		log.Printf("invalid password: %v", err)
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	log.Println("user exists")

	session, err := l.authS.NewSession(usr.Id)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	http.Redirect(w, r, "/chat", http.StatusSeeOther)

}
