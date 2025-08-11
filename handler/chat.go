package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
	"server/handler/socket"
	i "server/interfaces"
	"server/lib"
	cred "server/services/authService/credentials"
	typ "server/types"
)

type StatusCodeMessage = string

const (
	InternalServerError StatusCodeMessage = "Internal Server Error"
	msgConnectUserFail  StatusCodeMessage = "Unable to connect user"
	msgMalformedJSON    StatusCodeMessage = "Invalid JSON received"
	MethodNotAllowed    StatusCodeMessage = "request method not allowed"
)

type MessageType = string

const (
	NewMessage           MessageType = "1"
	EditMessage          MessageType = "2"
	DeleteMessage        MessageType = "3"
	NewContactMessage    MessageType = "4"
	EditContactMessage   MessageType = "5"
	DeleteContactMessage MessageType = "6"
)

type HTMLPath = string

const (
	chatViewHTML HTMLPath = "./static/templates/chat/chat-view.html"
	messagesHTML HTMLPath = "./static/templates/chat/messages.html"
	messageHTML  HTMLPath = "./static/templates/chat/message.html"
	chatHTML     HTMLPath = "./static/templates/chat/chat.html"
	chatsHTML    HTMLPath = "./static/templates/chat/chats.html"
	newChatHTML  HTMLPath = "./static/templates/chat/new-chat.html"
	contactHTML  HTMLPath = "./static/templates/chat/contact.html"
	contactsHTML HTMLPath = "./static/templates/chat/contacts.html"
)

var (
	chatViewTmpl *template.Template
)

func init() {
	chatViewTmpl = template.Must(
		template.ParseFiles(
			chatViewHTML,
			messagesHTML,
			messageHTML,
			chatHTML,
			chatsHTML,
			newChatHTML,
			contactHTML,
			contactsHTML,
		),
	)
}

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

func (h *ChatHandler) RenderChatPage(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
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
	userId := session.UserId()

	chats, err := h.chatS.GetChats(userId)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to get chats for userId: %v, Error: %v", userId, err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	var messages []entities.Message
	if len(chats) != 0 {
		chatId := chats[0].Id
		messages, err = h.msgS.GetChatMessages(chatId)
		if err != nil {
			h.lgr.LogError(fmt.Errorf("failed to get chat messages for chatId: %v, Error: %v", chatId, err))
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}
	}

	contacts, err := h.userS.GetContacts(userId)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to get contacts for userId: %v, Error: %v", userId, err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog(fmt.Sprintf("session userId: %v", userId))
	renderChatpayload := dto.RenderChatPayload{
		UserId:   userId,
		Chats:    chats,
		Messages: messages,
		Contacts: contacts,
	}

	if err = chatViewTmpl.ExecuteTemplate(w, "chatView", renderChatpayload); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to execute chatView template, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog("->>>> RESPONSE SENT")
}

/* MESSAGING ================================================================ */

func (h *ChatHandler) ChatWebsocket(w http.ResponseWriter, r *http.Request) {
	conn := socket.New(w, r)
	if conn == nil {
		h.lgr.LogError(errors.New("failed to connect user"))
		http.Error(w, msgConnectUserFail, http.StatusInternalServerError)
		return
	}

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		h.lgr.Log("user not authenticated, redirecting to landing page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		conn.Close()
		return
	}
	userId := session.UserId()

	h.connS.StoreConnection(conn, userId)
	defer h.connS.DisconnectUser(userId)

	h.lgr.DLog(fmt.Sprintf("active connections : %v", h.connS.GetActiveConnections()))

	for {

		payload, err := h.readIncomingMessage(conn)
		if err != nil {
			h.lgr.LogError(fmt.Errorf("failed to read incoming websocket message, Error: %v", err))
			http.Error(w, msgMalformedJSON, http.StatusBadRequest)
			break
		}

		switch payload.Type {
		case NewMessage:
			h.lgr.DLog("Handling new message...")

			newMessageRequest, err := payload.ParseNewMessageRequest()
			if err != nil {
				h.lgr.LogError(fmt.Errorf("Failed to parse message request %v", err))
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

			if err = h.handleNewMessageRequest(newMessageRequest, userId); err != nil {
				h.lgr.LogError(fmt.Errorf("Failed to handle message request %v", err))
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

		case NewContactMessage:
			h.lgr.DLog("Handling new contact message...")

			newMessageRequest, err := payload.ParseNewMessageRequest()
			if err != nil {
				h.lgr.LogError(fmt.Errorf("Failed to parse message request %v", err))
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

			if err = h.handleNewContactMessageRequest(newMessageRequest, userId); err != nil {
				h.lgr.LogError(fmt.Errorf("Failed to handle message request %v", err))
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				break
			}

		}
		h.lgr.DLog("->>>> RESPONSE SENT")

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

func (h *ChatHandler) handleNewMessageRequest(mr dto.NewMessageRequest, userId typ.UserId) error {
	h.lgr.LogFunctionInfo()

	h.lgr.DLog(fmt.Sprintf("chatId %v: replyId %v", mr.ChatId, mr.ReplyId))
	cid, err := lib.ConvertStringToInt64(mr.ChatId)
	if err != nil {
		return err
	}

	var rid int64
	if mr.ReplyId != "" {
		rid, err = lib.ConvertStringToInt64(mr.ReplyId)
		if err != nil {
			return err
		}
	}

	chatId := typ.ChatId(cid)
	replyId := typ.MessageId(rid)

	newMsgInput := dto.NewMessageInput{
		UserId:  userId,
		ChatId:  chatId,
		ReplyId: &replyId,
		Text:    mr.MsgText,
	}

	err = h.msgS.HandleNewMessage(newMsgInput)
	if err != nil {
		return err
	}

	return nil
}

func (h *ChatHandler) handleNewContactMessageRequest(mr dto.NewMessageRequest, userId typ.UserId) error {
	h.lgr.LogFunctionInfo()

	h.lgr.DLog(fmt.Sprintf("chatId %v: replyId %v", mr.ChatId, mr.ReplyId))
	cid, err := lib.ConvertStringToInt64(mr.ChatId)
	if err != nil {
		return err
	}

	var rid int64
	if mr.ReplyId != "" {
		rid, err = lib.ConvertStringToInt64(mr.ReplyId)
		if err != nil {
			return err
		}
	}

	chatId := typ.ChatId(cid)
	replyId := typ.MessageId(rid)

	newMsgInput := dto.NewMessageInput{
		UserId:  userId,
		ChatId:  chatId,
		ReplyId: &replyId,
		Text:    mr.MsgText,
	}

	err = h.msgS.HandleNewContactMessage(newMsgInput)
	if err != nil {
		return err
	}

	return nil
}

/* SWITCH CHAT ============================================================== */

func (h *ChatHandler) SwitchChat(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		h.lgr.LogError(errors.New("request method not allowed"))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	_, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		h.lgr.Log("user not authenticated, redirecting to landing page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	query := r.URL.Query()
	switchChatRequest := dto.SwitchChatRequest{
		ChatId: query.Get("chatid"),
	}

	switchChatResponse, err := h.handleChatSwitchRequest(switchChatRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle switch chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if err := chatViewTmpl.ExecuteTemplate(w, "messages", switchChatResponse); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to execute messages template for chat view, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog("->>>> RESPONSE SENT")
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
		Messages: messages,
	}

	return switchChatResponse, nil
}

/* NEW CHAT ================================================================= */

func (h *ChatHandler) NewChat(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query()
	newChatRequest := dto.NewChatRequest{
		Name: query.Get("name"),
	}

	newChatResponse, err := h.handleNewChatRequest(newChatRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle new chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if err := chatViewTmpl.ExecuteTemplate(w, "new-chat", newChatResponse); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to execute new-chat template for chat view, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleNewChatRequest(cr dto.NewChatRequest, userId typ.UserId) (dto.NewChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chat, err := h.chatS.NewChat(cr.Name, userId)
	if err != nil {
		return dto.NewChatResponse{}, err
	}

	newChatResponse := dto.NewChatResponse{
		Chats:    []entities.Chat{*chat},
		Messages: []entities.Message{},
	}

	return newChatResponse, nil
}

/* ADD FRIEND REQUEST ========================================================================== */

func (h *ChatHandler) AddContact(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query()
	addContactRequest := dto.AddContactRequest{
		Email: query.Get("email"),
	}

	addContactResponse, err := h.handleAddContactRequest(addContactRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle add contact request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if err := chatViewTmpl.ExecuteTemplate(w, "new-contact", addContactResponse); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to execute new-contact template for chat view, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog("->>>> RESPONSE SENT")
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

	conn := h.connS.GetConnection(typ.UserId(contact.Id))
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

func (h *ChatHandler) SwitchContactChat(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodPost {
		h.lgr.LogError(errors.New("request method not allowed"))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	_, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		h.lgr.Log("user not authenticated, redirecting to landing page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	query := r.URL.Query()
	switchChatRequest := dto.SwitchChatRequest{
		ChatId: query.Get("chatid"),
	}

	switchChatResponse, err := h.handleContactChatSwitchRequest(switchChatRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle contact chat switch request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if err := chatViewTmpl.ExecuteTemplate(w, "messages", switchChatResponse); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to execute messages template for chat view, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleContactChatSwitchRequest(switchChatRequest dto.SwitchChatRequest) (dto.SwitchChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := switchChatRequest.GetChatId()
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	messages, err := h.msgS.GetContactMessages(chatId)
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	switchChatResponse := dto.SwitchChatResponse{
		Messages: messages,
	}

	return switchChatResponse, nil
}
