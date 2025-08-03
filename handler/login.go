package handler

import (
	"encoding/json"
	"log"
	"net/http"
	dto "server/data/DTOs"
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

	var loginRequest dto.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}

	email, err := cred.NewEmail(loginRequest.Email)
	if err != nil {
		log.Printf("invalid email: %v", err)
		SendErrorResponse(w, "Invalid Email format", false)
		return
	}

	emails := []cred.Email{email}

	// Need to add some type of rollback in the case of an error?
	users, err := l.userS.FindUsers(emails)
	if err != nil {
		log.Printf("failed to find user %v", err)
		SendErrorResponse(w, "Email not found", false)
		return
	}

	if len(users) == 0 {
		log.Printf("failed to find user %v", err)
		SendErrorResponse(w, "Email not found", false)
		return
	}

	user := users[0]

	pwdBytes := []byte(loginRequest.Password)
	if err := user.PwdHash.Compare(pwdBytes); err != nil {
		log.Printf("invalid password: %v", err)
		SendErrorResponse(w, "Invalid password", false)
		return
	}

	session, err := l.authS.NewSession(user.Id)
	if err != nil {
		log.Printf("failed to create new user %v", err)
		http.Error(w, Status500, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	LoginSuccessful(w)
}

func LoginSuccessful(w http.ResponseWriter) {
	SendErrorResponse(w, "", true)
}

func SendErrorResponse(w http.ResponseWriter, message string, noError bool) {
	errorResponse := dto.ErrorResponse{
		NoError:      noError,
		ErrorMessage: message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse)
}
