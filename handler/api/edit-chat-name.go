package api

func EditChatName(a i.AuthService) http.Handler {
	h := editChatName{}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type editChatName struct {
}

func (h *ChatHandler) EditChatName(w http.ResponseWriter, r *http.Request) {
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

	var editChatNameRequest dto.ecnrequest
	if err := json.NewDecoder(r.Body).Decode(&editChatNameRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, ParseFormFail, http.StatusBadRequest)
		return
	}
	editChatNameRequest.UserId = session.UserId()

	editChatNameResponse, err := h.handleEditChatNameRequest(editChatNameRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle contact edit chat name request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, editChatNameResponse)

	h.lgr.DLog(fmt.Sprintf("->>>> RESPONSE SENT:: %v", editChatNameResponse))

}

func (h *ChatHandler) handleEditChatNameRequest(req dto.ecnrequest) (dto.ecnresponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := req.GetChatId()
	if err != nil {
		return dto.ecnresponse{}, err
	}

	err = h.chatS.EditChatName(req.Name, chatId, req.UserId)
	if err != nil {
		return dto.ecnresponse{}, err
	}

	return dto.ecnresponse{Name: req.Name}, nil
}

type ecnrequest struct {
	Name   string `json:"Name"`
	ChatId string `json:"ChatId"`
	UserId typ.UserId
}

type ecnresponse struct {
	Name string `json:"Name"`
}
