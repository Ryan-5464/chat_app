package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
	"server/handler/socket"
	i "server/interfaces"
	ss "server/services/authService/session"
	typ "server/types"
	"strconv"
	"text/template"

	ws "github.com/gorilla/websocket"
)

const (
	msgInternalServerError string = "Internal Server Error"
	msgConnectUserFail     string = "Unable to connect user"
	msgMalformedJSON       string = "Invalid JSON received"
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
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := r.Context().Value("session").(ss.Session)
	emptySession := ss.Session{}
	if session == emptySession {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	chats, err := cr.chatS.GetChats(session.UserId())
	if err != nil {
		log.Println("failed to get chat data for user")
		http.Error(w, "interrnal server error", http.StatusInternalServerError)
		return
	}

	var messages []entities.Message
	if len(chats) == 0 {
		messages = []entities.Message{}
	} else {
		messages, err = cr.msgS.GetChatMessages(chats[0].Id)
		if err != nil {
			log.Println("failed to get message data for user")
			http.Error(w, "interrnal server error", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("./static/templates/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Chats:    chats,
		Messages: messages,
	}

	log.Println("ONLOADCHAT:", data)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cr *ChatHandler) ChatWebsocket(w http.ResponseWriter, r *http.Request) {
	conn := socket.New(w, r)
	if conn == nil {
		http.Error(w, msgConnectUserFail, http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	session, userAuthenticated := checkAuthenticationStatus(r)
	if !userAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cr.connS.StoreConnection(conn, session.UserId())

	for {

		payload, err := cr.readIncomingMessage(conn)
		if err != nil {
			http.Error(w, msgMalformedJSON, http.StatusBadRequest)
			break
		}

		switch payload.Type {
		case "SwitchChat":
			log.Println("Switching chat")

			msgPayload, err := cr.HandleChatSwitch(payload.Data)
			if err != nil {
				log.Println(err)
				http.Error(w, msgInternalServerError, http.StatusInternalServerError)
				break
			}

			err = conn.WriteMessage(ws.TextMessage, msgPayload)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

		case "NewMessage":
			log.Println("NewMessage")

			newMsg, err := cr.parseNewMessageData(payload.Data)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			log.Println("newMsg", newMsg)

			if cr.msgS == nil {
				log.Fatal("missing messaging service")
			}

			if err := cr.msgS.HandleNewMessage(newMsg); err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				break
			}

		case "NewChat":

			msgPayload, err := cr.newChat(session.UserId(), payload.Data)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			err = conn.WriteMessage(ws.TextMessage, msgPayload)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

		}
	}
	log.Println("websocket closed")

}

func (cr *ChatHandler) readIncomingMessage(conn i.Socket) (dto.Payload, error) {
	cr.lgr.DLog("Reading incoming message...")
	payload := dto.Payload{}

	err := conn.ReadJSON(&payload)
	if err != nil {
		errReadJSONFail := fmt.Errorf("Failed to read JSON: %v", err)
		cr.lgr.LogError(errReadJSONFail)
		return payload, err
	}

	msgType := fmt.Sprintf("Message type: %v", payload.Type)
	cr.lgr.DLog(msgType)
	return payload, nil
}

func parseChatSwitchRequest(data []byte) (dto.SwitchChat, error) {
	s := dto.SwitchChat{}
	if err := json.Unmarshal(data, &s); err != nil {
		return s, err
	}
	return s, nil
}

func parseNewChatData(data []byte) (entities.Chat, error) {
	n := dto.NewChat{}
	if err := json.Unmarshal(data, &n); err != nil {
		return entities.Chat{}, err
	}

	chat := entities.Chat{
		Name: n.Name,
	}

	return chat, nil
}

func (cr *ChatHandler) parseNewMessageData(data []byte) (entities.Message, error) {
	cr.lgr.LogFunctionInfo()

	n := dto.NewMessage{}
	if err := json.Unmarshal(data, &n); err != nil {
		errUnmarshal := fmt.Errorf("failed to unmarshal data: %w", err)
		cr.lgr.LogError(errUnmarshal)
		return entities.Message{}, err
	}

	userId, err := convertStringToInt64(n.UserId)
	if err != nil {
		errTypeConversion := fmt.Errorf("failed to convert userId from string to int: %w", err)
		cr.lgr.LogError(errTypeConversion)
		return entities.Message{}, err
	}

	chatId, err := convertStringToInt64(n.ChatId)
	if err != nil {
		errTypeConversion := fmt.Errorf("failed to convert chatId from string to int: %w", err)
		cr.lgr.LogError(errTypeConversion)
		return entities.Message{}, err
	}

	var replyId int64
	if n.ReplyId != "" {
		replyId, err = convertStringToInt64(n.ReplyId)
		if err != nil {
			errTypeConversion := fmt.Errorf("failed to convert replyId from string to int: %w", err)
			cr.lgr.LogError(errTypeConversion)
			return entities.Message{}, err
		}
	}

	msgE := entities.Message{
		UserId:  typ.UserId(userId),
		ChatId:  typ.ChatId(chatId),
		ReplyId: typ.MessageId(replyId),
		Text:    n.MsgText,
	}

	return msgE, nil
}

func convertStringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func (cr *ChatHandler) HandleChatSwitch(switchChatrequest []byte) ([]byte, error) {
	cr.lgr.LogFunctionInfo()

	data, err := parseChatSwitchRequest(switchChatrequest)
	if err != nil {
		return []byte{}, err
	}

	chatId, err := data.GetChatId()
	if err != nil {
		return []byte{}, err
	}

	messages, err := cr.msgS.GetChatMessages(chatId)
	if err != nil {
		return []byte{}, err
	}

	payload := struct {
		Type     string
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Type:     "SwitchChat",
		Chats:    nil,
		Messages: messages,
	}

	msgPayload, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return msgPayload, nil
}

func (cr *ChatHandler) newChat(userId typ.UserId, newChatData []byte) ([]byte, error) {
	cr.lgr.LogFunctionInfo()

	newChat, err := parseNewChatData(newChatData)
	if err != nil {
		return []byte{}, err
	}

	newChat.AdminId = userId

	chat, err := cr.chatS.NewChat(newChat)
	if err != nil {
		return []byte{}, err
	}

	newMsg := entities.Message{
		ChatId: chat.Id,
		UserId: chat.AdminId,
		Text:   "Chat Created",
	}

	msg, err := cr.msgS.NewMessage(newMsg)
	if err != nil {
		return []byte{}, err
	}

	payload := struct {
		Type     string
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Type:     "NewChat",
		Chats:    []entities.Chat{chat},
		Messages: []entities.Message{msg},
	}

	msgPayload, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	return msgPayload, nil
}
