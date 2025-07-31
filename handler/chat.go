package handler

import (
	"encoding/json"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
	"server/handler/socket"
	i "server/interfaces"
	typ "server/types"
	"text/template"
)

const (
	InternalServerError string = "Internal Server Error"
	msgConnectUserFail  string = "Unable to connect user"
	msgMalformedJSON    string = "Invalid JSON received"
	MethodNotAllowed    string = "request method not allowed"
)

type MessageType = string

const (
	NewMessage    MessageType = "1"
	EditMessage   MessageType = "2"
	DeleteMessage MessageType = "3"
)

func NewChatHandler(lgr i.Logger, a i.AuthService, c i.ChatService, m i.MessageService, cnx i.ConnectionService, u i.UserService) *ChatHandler {
	return &ChatHandler{
		lgr:   lgr,
		authS: a,
		chatS: c,
		msgS:  m,
		connS: cnx,
		userS: u,
	}
}

type ChatHandler struct {
	lgr   i.Logger
	authS i.AuthService
	chatS i.ChatService
	msgS  i.MessageService
	connS i.ConnectionService
	userS i.UserService
}

func (cr *ChatHandler) RenderChatPage(w http.ResponseWriter, r *http.Request) {
	cr.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	chats, err := cr.chatS.GetChats(session.UserId())
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	var messages []entities.Message
	if len(chats) != 0 {
		messages, err = cr.msgS.GetChatMessages(chats[0].Id)
		if err != nil {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("./static/templates/chat.html")
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	renderChatpayload := dto.RenderChatPayload{
		Chats:    chats,
		Messages: messages,
	}

	err = tmpl.Execute(w, renderChatpayload)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}
}

/* MESSAGING ================================================================ */

func (h *ChatHandler) ChatWebsocket(w http.ResponseWriter, r *http.Request) {
	conn := socket.New(w, r)
	if conn == nil {
		http.Error(w, msgConnectUserFail, http.StatusInternalServerError)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		conn.Close()
		return
	}

	h.connS.StoreConnection(conn, session.UserId())
	defer h.connS.DisconnectUser(session.UserId())

	for {

		payload, err := h.readIncomingMessage(conn)
		if err != nil {
			http.Error(w, msgMalformedJSON, http.StatusBadRequest)
			break
		}

		switch payload.Type {
		case NewMessage:
			h.lgr.DLog("Handling new message...")

			newMessageRequest, err := payload.ParseNewMessageRequest()
			if err != nil {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

			if err = h.handleNewMessageRequest(newMessageRequest); err != nil {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

		}

	}

	h.lgr.Log("User disconnected")
}

func (h *ChatHandler) readIncomingMessage(conn i.Socket) (dto.WebsocketPayload, error) {
	h.lgr.LogFunctionInfo()

	payload := dto.WebsocketPayload{}
	err := conn.ReadJSON(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (h *ChatHandler) handleNewMessageRequest(newMessageRequest dto.NewMessageRequest) error {
	msgE, err := newMessageRequest.ToMessageEntity()
	if err != nil {
		return err
	}

	err = h.msgS.HandleNewMessage(msgE)
	if err != nil {
		return err
	}

	return nil
}

/* SWITCH CHAT ============================================================== */

func (h *ChatHandler) SwitchChat(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	_, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var switchChatRequest dto.SwitchChatRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&switchChatRequest); err != nil {
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}

	switchChatResponse, err := h.handleChatSwitchRequest(switchChatRequest)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(switchChatResponse)
}

func (h *ChatHandler) handleChatSwitchRequest(switchChatRequest dto.SwitchChatRequest) (dto.SwitchChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := switchChatRequest.GetChatId()
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	messages, err := h.msgS.GetChatMessages(chatId)
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	switchChatResponse := dto.SwitchChatResponse{
		Messages:        messages,
		NewActiveChatId: chatId,
	}

	return switchChatResponse, nil
}

/* NEW CHAT ================================================================= */

func (h *ChatHandler) NewChat(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var newChatRequest dto.NewChatRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newChatRequest); err != nil {
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}

	newChatResponse, err := h.handleNewChatRequest(newChatRequest, session.UserId())
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newChatResponse)
}

func (cr *ChatHandler) handleNewChatRequest(newChatRequest dto.NewChatRequest, userId typ.UserId) (dto.NewChatResponse, error) {
	cr.lgr.LogFunctionInfo()

	newChat := entities.Chat{
		Name: newChatRequest.Name,
	}
	newChat.AdminId = userId

	chat, err := cr.chatS.NewChat(newChat)
	if err != nil {
		return dto.NewChatResponse{}, err
	}

	newChatResponse := dto.NewChatResponse{
		Chats:           []entities.Chat{chat},
		Messages:        []entities.Message{},
		NewActiveChatId: chat.Id,
	}

	return newChatResponse, nil
}
