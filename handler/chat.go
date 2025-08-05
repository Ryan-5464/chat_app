package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
	"server/handler/socket"
	i "server/interfaces"
	cred "server/services/authService/credentials"
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

	contacts, err := cr.userS.GetContacts(session.UserId())
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./static/templates/chat.html")
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	renderChatpayload := dto.RenderChatPayload{
		Chats:    chats,
		Messages: messages,
		Contacts: contacts,
	}

	if err = tmpl.Execute(w, renderChatpayload); err != nil {
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
	userId := session.UserId()

	h.connS.StoreConnection(conn, userId)
	defer h.connS.DisconnectUser(userId)

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

			if err = h.handleNewMessageRequest(newMessageRequest, userId); err != nil {
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

func (h *ChatHandler) handleNewMessageRequest(newMessageRequest dto.NewMessageRequest, userId typ.UserId) error {
	msgE, err := newMessageRequest.ToMessageEntity(userId)
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
		log.Println("MALFORMEDJSON:::::")
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
	// move this check close to the database
	if messages == nil {
		messages = []entities.Message{}
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

func (h *ChatHandler) handleNewChatRequest(newChatRequest dto.NewChatRequest, userId typ.UserId) (dto.NewChatResponse, error) {
	h.lgr.LogFunctionInfo()

	newChat := entities.Chat{
		Name: newChatRequest.Name,
	}
	newChat.AdminId = userId

	chat, err := h.chatS.NewChat(newChat)
	if err != nil {
		return dto.NewChatResponse{}, err
	}

	newChatResponse := dto.NewChatResponse{
		Chats:           []entities.Chat{*chat},
		Messages:        []entities.Message{},
		NewActiveChatId: chat.Id,
	}

	return newChatResponse, nil
}

/* ADD FRIEND REQUEST ========================================================================== */

func (h *ChatHandler) AddContact(w http.ResponseWriter, r *http.Request) {
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

	var addContactRequest dto.AddContactRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&addContactRequest); err != nil {
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}

	addContactResponse, err := h.handleAddContactRequest(addContactRequest, session.UserId())
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addContactResponse)
}

func (h *ChatHandler) handleAddContactRequest(addContactRequest dto.AddContactRequest, userId typ.UserId) (dto.AddContactResponse, error) {
	h.lgr.LogFunctionInfo()

	var addContactResponse dto.AddContactResponse

	addContactInput := dto.AddContactInput{
		Email:  cred.Email(addContactRequest.Email),
		UserId: userId,
	}

	contact, err := h.userS.AddContact(addContactInput)
	if err != nil {
		return addContactResponse, err
	}

	if contact == nil {
		return addContactResponse, errors.New("failed to add contact")
	}

	conn := h.connS.GetConnection(contact.Id)
	if conn == nil {
		contact.OnlineStatus = false
	} else {
		contact.OnlineStatus = true
	}

	addContactResponse = dto.AddContactResponse{
		Contacts: []entities.Contact{*contact},
	}

	return addContactResponse, nil
}

func (h *ChatHandler) SwitchPrivateChat(w http.ResponseWriter, r *http.Request) {
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
		log.Println("MALFORMEDJSON:::::")
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}

	switchChatResponse, err := h.handlePrivateChatSwitchRequest(switchChatRequest)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(switchChatResponse)
}

func (h *ChatHandler) handlePrivateChatSwitchRequest(switchChatRequest dto.SwitchChatRequest) (dto.SwitchChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := switchChatRequest.GetChatId()
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	messages, err := h.msgS.GetChatMessages(chatId)
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}
	// move this check close to the database
	if messages == nil {
		messages = []entities.Message{}
	}

	switchChatResponse := dto.SwitchChatResponse{
		Messages:        messages,
		NewActiveChatId: chatId,
	}

	return switchChatResponse, nil
}
