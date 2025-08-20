package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	dto "server/data/DTOs"
	i "server/interfaces"
	typ "server/types"
)

func NewProfileHandler(l i.Logger, a i.AuthService, u i.UserService) *ProfileHandler {
	return &ProfileHandler{
		lgr:   l,
		authS: a,
		userS: u,
	}
}

type ProfileHandler struct {
	authS i.AuthService
	userS i.UserService
	lgr   i.Logger
}

func (h *ProfileHandler) RenderProfilePage(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()
	log.Println("RENDER PROFILE PAGE ====> ")

	if r.Method != http.MethodGet {
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		h.lgr.Log("user not authenticated, redirecting to landing page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println(fmt.Sprintf("USERID ====> %v", session.UserId()))

	user, err := h.userS.GetUser(session.UserId())
	if err != nil {
		h.lgr.Log(fmt.Sprintf("Failed to get user for user id %v", session.UserId()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(fmt.Sprintf("USER ====> %v", user))

	log.Println(fmt.Sprintf("USERNAME ====> %v", user.Name))

	data := struct {
		Name string
	}{
		Name: user.Name,
	}

	tmpl, err := template.ParseFiles("./static/pages/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProfileHandler) EditUserName(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		h.lgr.LogError(errors.New("request method not allowed"))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		h.lgr.Log("user not authenticated, redirecting to landing page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var editUserNameRequest dto.EditUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&editUserNameRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, ParseFormFail, http.StatusBadRequest)
		return
	}
	userId := session.UserId()

	editUserNameResponse, err := h.handleEditUserNameRequest(editUserNameRequest, userId)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle contact edit chat name request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, editUserNameResponse)

	h.lgr.DLog(fmt.Sprintf("->>>> RESPONSE SENT:: %v", editUserNameResponse))

}

func (h *ProfileHandler) handleEditUserNameRequest(req dto.EditUserNameRequest, userId typ.UserId) (dto.EditUserNameResponse, error) {
	h.lgr.LogFunctionInfo()

	err := h.userS.EditUserName(req.Name, userId)
	if err != nil {
		return dto.EditUserNameResponse{}, err
	}

	return dto.EditUserNameResponse{Name: req.Name}, nil
}
