package renderers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/data/entities"
	i "server/interfaces"
	typ "server/types"
	"text/template"

	ws "github.com/gorilla/websocket"
)

func NewChatRenderer(a i.AuthService, c i.ChatService, m i.MessageService) *ChatRenderer {
	return &ChatRenderer{
		authS: a,
		chatS: c,
		msgS:  m,
	}
}

type ChatRenderer struct {
	authS i.AuthService
	chatS i.ChatService
	msgS  i.MessageService
}

func (cr *ChatRenderer) RenderChat(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("session_token")
	// if err != nil {
	// 	log.Println("no session cookie found")
	// 	http.Error(w, "no session cookie found", http.StatusInternalServerError)
	// 	return
	// }
	testToken, err := testToken(cr.authS)
	if err != nil {
		log.Println("failed to create dummy session")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	session, err := cr.authS.ValidateAndRefreshSession(testToken)
	if err != nil {
		log.Println("failed to valdiate session token")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, session.Cookie())

	chats, err := cr.chatS.GetChats()
	if err != nil {
		log.Println("failed to get chat data for user")
		http.Error(w, "interrnal server error", http.StatusInternalServerError)
		return
	}

	messages, err := cr.msgS.GetMessages(typ.ChatId(1))
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

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cr *ChatRenderer) ChatWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := establishWebsocket(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer conn.Close()

	for {
		pl := payload{}
		if err := conn.ReadJSON(&pl); err != nil {
			log.Println("failed to read JSON: ", err)
			http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
			break
		}

		log.Println("pl", pl)

		switch pl.Type {
		case "SwitchChat":
			log.Println("Switching chat")

			chat := entities.Chat{}
			if err := json.Unmarshal(pl.Data, &chat); err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			log.Println("chat", chat)

			messages, err := cr.msgS.GetMessages(typ.ChatId(chat.Id))
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}

			log.Println("messages", messages)

			data := struct {
				Chats    []entities.Chat
				Messages []entities.Message
			}{
				Chats:    nil,
				Messages: messages,
			}

			msgPayload, err := json.Marshal(data)
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
		log.Println("websocket closed")
	}
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
