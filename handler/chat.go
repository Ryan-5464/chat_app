package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
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
	memberModalHTML     HTMLPath = "./static/templates/chat/member-modal.html"
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
			memberModalHTML,
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
	h.lgr.Log(fmt.Sprintf("Add member email: %v", addMemberToChatRequest.Email))
	h.lgr.Log(fmt.Sprintf("Chat id: %v", addMemberToChatRequest.ChatId))

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

	h.lgr.DLog(fmt.Sprintf("NEW MEMBER ID => ", memberId))

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
