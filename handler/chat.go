package handler

import (
	"encoding/json"
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
	ParseFormFail       StatusCodeMessage = "Failed to parse form"
	UnauthorizedRequest StatusCodeMessage = "User not authorized to make change"
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
	chatViewHTML        HTMLPath = "./static/templates/chat/chat-view.html"
	messagesHTML        HTMLPath = "./static/templates/chat/messages.html"
	messageHTML         HTMLPath = "./static/templates/chat/message.html"
	chatHTML            HTMLPath = "./static/templates/chat/chat.html"
	chatsHTML           HTMLPath = "./static/templates/chat/chats.html"
	newChatHTML         HTMLPath = "./static/templates/chat/new-chat.html"
	newChatNameHTML     HTMLPath = "./static/templates/chat/new-chat-name.html"
	LeaveChatHTML       HTMLPath = "./static/templates/chat/leave-chat.html"
	contactHTML         HTMLPath = "./static/templates/chat/contact.html"
	contactsHTML        HTMLPath = "./static/templates/chat/contacts.html"
	contactModalHTML    HTMLPath = "./static/templates/chat/contact-modal.html"
	chatModalHTML       HTMLPath = "./static/templates/chat/chat-modal.html"
	messageModalHTML    HTMLPath = "./static/templates/chat/message-modal.html"
	memberListModalHTML HTMLPath = "./static/templates/chat/member-list-modal.html"
)

var (
	chatViewTmpl *template.Template
)

func init() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": dict,
	})
	chatViewTmpl = template.Must(
		tmpl.ParseFiles(
			chatViewHTML,
			messagesHTML,
			messageHTML,
			chatHTML,
			chatsHTML,
			newChatHTML,
			newChatNameHTML,
			LeaveChatHTML,
			contactHTML,
			contactsHTML,
			contactModalHTML,
			chatModalHTML,
			messageModalHTML,
			memberListModalHTML,
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

	var chatId typ.ChatId = -1
	var messages []entities.Message
	if len(chats) != 0 {
		chatId = chats[0].Id
		messages, err = h.msgS.GetChatMessages(chatId, userId)
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
		UserId:       userId,
		Chats:        chats,
		Messages:     messages,
		Contacts:     contacts,
		ActiveChatId: chatId,
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

// GET CHAT MEMBERS ============================================================

func (h *ChatHandler) GetChatMembers(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
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
	getChatmembersRequest := dto.GetChatMembersRequest{
		ChatId: query.Get("ChatId"),
	}

	getChatMembersRequest, err := h.handleGetChatMembersRequest(getChatmembersRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle switch chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, getChatMembersRequest)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleGetChatMembersRequest(mr dto.GetChatMembersRequest) (dto.GetChatMembersResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := lib.ConvertStringToInt64(mr.ChatId)
	if err != nil {
		return dto.GetChatMembersResponse{}, err
	}

	members, err := h.chatS.GetChatMembers(typ.ChatId(chatId))
	if err != nil {
		return dto.GetChatMembersResponse{}, err
	}

	getChatMembersResponse := dto.GetChatMembersResponse{
		Members: members,
	}

	return getChatMembersResponse, nil
}

/* SWITCH CHAT ============================================================== */

func (h *ChatHandler) SwitchChat(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query()
	switchChatRequest := dto.SwitchChatRequest{
		ChatId: query.Get("ChatId"),
	}

	switchChatResponse, err := h.handleChatSwitchRequest(switchChatRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle switch chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, switchChatResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleChatSwitchRequest(switchChatRequest dto.SwitchChatRequest, userId typ.UserId) (dto.SwitchChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := switchChatRequest.GetChatId()
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	messages, err := h.msgS.GetChatMessages(chatId, userId)
	if err != nil {
		return dto.SwitchChatResponse{}, err
	}

	switchChatResponse := dto.SwitchChatResponse{
		ActiveChatId: chatId,
		Messages:     messages,
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

	var newChatRequest dto.NewChatRequest
	if err := json.NewDecoder(r.Body).Decode(&newChatRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}
	h.lgr.DLog(fmt.Sprintf("New chat request name: %v", newChatRequest.Name))

	newChatResponse, err := h.handleNewChatRequest(newChatRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle new chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, newChatResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleNewChatRequest(cr dto.NewChatRequest, userId typ.UserId) (dto.NewChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chat, err := h.chatS.NewChat(cr.Name, userId)
	if err != nil {
		return dto.NewChatResponse{}, err
	}

	newChatResponse := dto.NewChatResponse{
		Chats: []entities.Chat{*chat},
	}

	return newChatResponse, nil
}

/* REMOVE CONTACT REQUEST ====================================================================== */

func (h *ChatHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodDelete {
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
	removeContactRequest := dto.RemoveContactRequest{
		ContactId: query.Get("ContactId"),
	}

	removeContactResponse, err := h.handleRemoveContactRequest(removeContactRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle remove contact request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, removeContactResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleRemoveContactRequest(cr dto.RemoveContactRequest, userId typ.UserId) (dto.RemoveContactResponse, error) {
	h.lgr.LogFunctionInfo()

	contactId, err := lib.ConvertStringToInt64(cr.ContactId)
	if err != nil {
		return dto.RemoveContactResponse{}, err
	}

	if err := h.userS.RemoveContact(typ.ContactId(contactId), userId); err != nil {
		return dto.RemoveContactResponse{}, err
	}

	chats, err := h.chatS.GetChats(userId)
	if err != nil {
		return dto.RemoveContactResponse{}, err
	}

	newActiveChatId := chats[0].Id
	var messages []entities.Message
	if len(chats) != 0 {
		chatId := newActiveChatId
		messages, err = h.msgS.GetChatMessages(chatId, userId)
		if err != nil {
			return dto.RemoveContactResponse{}, err
		}
	}

	contacts, err := h.userS.GetContacts(userId)
	if err != nil {
		return dto.RemoveContactResponse{}, err
	}

	newChatResponse := dto.RemoveContactResponse{
		NewActiveChatId: newActiveChatId,
		Contacts:        contacts,
		Messages:        messages,
	}

	return newChatResponse, nil
}

/* LEAVE CHAT REQUEST ========================================================================== */

func (h *ChatHandler) LeaveChat(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodDelete {
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
	leaveChatRequest := dto.LeaveChatRequest{
		ChatId: query.Get("ChatId"),
	}

	leaveChatResponse, err := h.handleLeaveChatRequest(leaveChatRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle leave chat request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, leaveChatResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleLeaveChatRequest(cr dto.LeaveChatRequest, userId typ.UserId) (dto.LeaveChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := lib.ConvertStringToInt64(cr.ChatId)
	if err != nil {
		return dto.LeaveChatResponse{}, err
	}

	chats, err := h.chatS.LeaveChat(typ.ChatId(chatId), userId)
	if err != nil {
		return dto.LeaveChatResponse{}, err
	}

	newActiveChatId := chats[0].Id
	var messages []entities.Message
	if len(chats) != 0 {
		chatId := newActiveChatId
		messages, err = h.msgS.GetChatMessages(chatId, userId)
		if err != nil {
			return dto.LeaveChatResponse{}, err
		}
	}

	newChatResponse := dto.LeaveChatResponse{
		NewActiveChatId: newActiveChatId,
		Chats:           chats,
		Messages:        messages,
	}

	return newChatResponse, nil
}

// ADD MEMBER TO CHAT ============================================================================

func (h *ChatHandler) AddMemberToChat(w http.ResponseWriter, r *http.Request) {
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

	var addMemberToChatRequest dto.AddMemberToChatRequest
	if err := json.NewDecoder(r.Body).Decode(&addMemberToChatRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}
	h.lgr.DLog(fmt.Sprintf("Add member email: %v", addMemberToChatRequest.Email))
	h.lgr.DLog(fmt.Sprintf("Chat id: %v", addMemberToChatRequest.ChatId))

	addMemberToChatResponse, err := h.handleAddMemberToChatRequest(addMemberToChatRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle add member request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, addMemberToChatResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")

}

func (h *ChatHandler) handleAddMemberToChatRequest(ar dto.AddMemberToChatRequest) (dto.AddMemberToChatResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := lib.ConvertStringToInt64(ar.ChatId)
	if err != nil {
		return dto.AddMemberToChatResponse{}, err
	}

	memberId, err := h.chatS.AddMember(cred.Email(ar.Email), typ.ChatId(chatId))
	if err != nil {
		return dto.AddMemberToChatResponse{}, err
	}

	member, err := h.chatS.GetChatMember(typ.ChatId(chatId), memberId)
	if err != nil {
		return dto.AddMemberToChatResponse{}, err
	}

	addMemberToChatResponse := dto.AddMemberToChatResponse{
		Members: []entities.Member{*member},
	}

	return addMemberToChatResponse, nil
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

	var addContactRequest dto.AddContactRequest
	if err := json.NewDecoder(r.Body).Decode(&addContactRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, msgMalformedJSON, http.StatusBadRequest)
		return
	}
	h.lgr.DLog(fmt.Sprintf("Add contact request name: %v", addContactRequest.Email))

	addContactResponse, err := h.handleAddContactRequest(addContactRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle add contact request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, addContactResponse)

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

	query := r.URL.Query()
	switchContactChatRequest := dto.SwitchContactChatRequest{
		ContactChatId: query.Get("ContactChatId"),
	}

	switchContactChatResponse, err := h.handleContactChatSwitchRequest(switchContactChatRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle contact chat switch request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, switchContactChatResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleContactChatSwitchRequest(s dto.SwitchContactChatRequest, userId typ.UserId) (dto.SwitchContactChatResponse, error) {
	h.lgr.LogFunctionInfo()

	cid, err := lib.ConvertStringToInt64(s.ContactChatId)
	if err != nil {
		return dto.SwitchContactChatResponse{}, err
	}
	contactChatId := typ.ChatId(cid)

	messages, err := h.msgS.GetContactMessages(contactChatId, userId)
	if err != nil {
		return dto.SwitchContactChatResponse{}, err
	}

	switchChatResponse := dto.SwitchContactChatResponse{
		ActiveContactChatId: contactChatId,
		Messages:            messages,
	}

	return switchChatResponse, nil
}

// EDIT MESSAGE ====================================================================

func (h *ChatHandler) EditMessage(w http.ResponseWriter, r *http.Request) {
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

	var editMessageRequest dto.EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&editMessageRequest); err != nil {
		h.lgr.LogError(fmt.Errorf("failed to decode JSON request body: ", err))
		http.Error(w, ParseFormFail, http.StatusBadRequest)
		return
	}

	userId, err := lib.ConvertStringToInt64(editMessageRequest.UserId)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.UserId() != typ.UserId(userId) {
		h.lgr.LogError(fmt.Errorf("user does not own message: ", err))
		http.Error(w, UnauthorizedRequest, http.StatusBadRequest)
		return
	}

	editMessageResponse, err := h.handleEditMessageRequest(editMessageRequest)
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to handle contact edit chat name request, Error: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, editMessageResponse)

	h.lgr.DLog(fmt.Sprintf("->>>> RESPONSE SENT:: %v", editMessageResponse))

}

func (h *ChatHandler) handleEditMessageRequest(mr dto.EditMessageRequest) (dto.EditMessageResponse, error) {
	h.lgr.LogFunctionInfo()

	messageId, err := lib.ConvertStringToInt64(mr.MessageId)
	if err != nil {
		return dto.EditMessageResponse{}, err
	}

	msg, err := h.msgS.EditMessage(mr.MsgText, typ.MessageId(messageId))
	if err != nil {
		return dto.EditMessageResponse{}, err
	}

	return dto.EditMessageResponse{
		MsgText: msg.Text,
	}, nil

}

// EDIT CHAT NAME ==================================================================

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

	var editChatNameRequest dto.EditChatNameRequest
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

func (h *ChatHandler) handleEditChatNameRequest(req dto.EditChatNameRequest) (dto.EditChatNameResponse, error) {
	h.lgr.LogFunctionInfo()

	chatId, err := req.GetChatId()
	if err != nil {
		return dto.EditChatNameResponse{}, err
	}

	err = h.chatS.EditChatName(req.Name, chatId, req.UserId)
	if err != nil {
		return dto.EditChatNameResponse{}, err
	}

	return dto.EditChatNameResponse{Name: req.Name}, nil
}

func (h *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	h.lgr.LogFunctionInfo()

	if r.Method != http.MethodDelete {
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
	userId, err := lib.ConvertStringToInt64(query.Get("UserId"))
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to parse userId: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.UserId() != typ.UserId(userId) {
		h.lgr.Log("user unauthorized to delete message")
		return
	}

	deleteMessageRequest := dto.DeleteMessageRequest{
		MessageId: query.Get("MessageId"),
		UserId:    query.Get("UserId"),
	}

	deleteMessageResponse, err := h.handleDeleteMessageRequest(deleteMessageRequest, session.UserId())
	if err != nil {
		h.lgr.LogError(fmt.Errorf("failed to delete message: %v", err))
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, deleteMessageResponse)

	h.lgr.DLog("->>>> RESPONSE SENT")
}

func (h *ChatHandler) handleDeleteMessageRequest(dr dto.DeleteMessageRequest, userId typ.UserId) (dto.DeleteMessageResponse, error) {
	h.lgr.LogFunctionInfo()

	messageId, err := lib.ConvertStringToInt64(dr.MessageId)
	if err != nil {
		return dto.DeleteMessageResponse{}, err
	}

	if err := h.msgS.DeleteMessage(typ.MessageId(messageId)); err != nil {
		return dto.DeleteMessageResponse{}, err
	}

	return dto.DeleteMessageResponse{
		MessageId: typ.MessageId(messageId),
	}, nil
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: must have even number of args")
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}
