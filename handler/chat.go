package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	dto "server/data/DTOs"
	"server/data/entities"
	i "server/interfaces"
	typ "server/types"
	"strconv"
	"text/template"

	ws "github.com/gorilla/websocket"
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

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if cookie == nil {
		log.Println("No cookie found, redirecting to index page...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	token := cookie.Value
	session, err := cr.authS.ValidateAndRefreshSession(token)
	if err != nil {
		log.Println("error validating or refreshing session", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	chats, err := cr.chatS.GetChats(session.UserId())
	if err != nil {
		log.Println("failed to get chat data for user")
		http.Error(w, "interrnal server error", http.StatusInternalServerError)
		return
	}

	messages, err := cr.msgS.GetChatMessages(typ.ChatId(1))
	if err != nil {
		log.Println("failed to get message data for user")
		http.Error(w, "interrnal server error", http.StatusInternalServerError)
		return
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
	conn, err := establishWebsocket(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer conn.Close()

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	session, err := cr.authS.ValidateAndRefreshSession(cookie.Value)
	if err != nil {
		log.Println("failed to valdiate session token")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	cr.connS.StoreConnection(conn, session.UserId())

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read JSON: ", err)
			http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
			break
		}

		pl, err := parsePayload(payload)
		if err != nil {
			log.Println("failed to read JSON: ", err)
			http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
			break
		}

		log.Println("pl", pl)

		switch pl.Type {
		case "SwitchChat":
			log.Println("Switching chat")

			data, err := parseSwitchChatData(pl.Data)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			log.Println("chat", data)

			chatId, err := convertStringToInt64(data.ChatId)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			messages, err := cr.msgS.GetChatMessages(typ.ChatId(chatId))
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			log.Println("messages", messages)

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

		case "NewMessage":
			log.Println("NewMessage")

			newMsg, err := cr.parseNewMessageData(pl.Data)
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

			newChat, err := parseNewChatData(pl.Data)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			newChat.AdminId = session.UserId()

			cr.lgr.Log(fmt.Sprintf("%v", newChat))

			chat, err := cr.chatS.NewChat(newChat)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to create new chat", http.StatusInternalServerError)
				break
			}

			newMsg := entities.Message{
				ChatId: chat.Id,
				UserId: chat.AdminId,
				Text:   "Hello",
			}

			msg, err := cr.msgS.NewMessage(newMsg)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to create new message", http.StatusInternalServerError)
				break
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

func establishWebsocket(w http.ResponseWriter, r *http.Request) (*ws.Conn, error) {
	var upgrader = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to establish websocket connection %w", err)
	}

	return conn, nil
}

func testToken(authS i.AuthService) (string, error) {
	testSession, err := authS.NewSession(typ.UserId(1))
	if err != nil {
		return "", fmt.Errorf("failed to create test session : %w", err)
	}
	return testSession.JWEToken(), nil
}

func parsePayload(payload []byte) (dto.Payload, error) {
	p := dto.Payload{}
	if err := json.Unmarshal(payload, &p); err != nil {
		return p, err
	}
	return p, nil
}

func parseSwitchChatData(data []byte) (dto.SwitchChat, error) {
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
