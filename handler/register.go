package handler

import (
	"log"
	"net/http"
	dto "server/data/DTOs"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	ss "server/services/authService/session"
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

	session := r.Context().Value("session").(ss.Session)
	emptySession := ss.Session{}
	if session != emptySession {
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

	newUserInput := dto.NewUserInput{
		Email:   email,
		PwdHash: pwdHash,
	}

	// Need to add some type of rollback in the case of an error
	user, err := rh.userS.NewUser(newUserInput)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("successfully created new user")

	session, err := rh.authS.NewSession(user.Id)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
