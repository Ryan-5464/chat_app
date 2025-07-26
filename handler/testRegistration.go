package handler

import (
	"log"
	"net/http"
	td "server/data/test"
	i "server/interfaces"
)

func NewTestRegistrationHandler(l i.Logger, u i.UserService) *TestRegistrationHandler {
	return &TestRegistrationHandler{
		lgr:   l,
		userS: u,
	}
}

type TestRegistrationHandler struct {
	lgr   i.Logger
	userS i.UserService
}

func (rh *TestRegistrationHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	rh.lgr.LogFunctionInfo()

	usr := td.TestUser()
	newUser, err := rh.userS.NewUser(usr)
	if err != nil {
		log.Println("User registration failed", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("new user", newUser)

	// http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
