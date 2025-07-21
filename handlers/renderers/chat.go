package renderers

import (
	"fmt"
	"log"
	"net/http"
	"server/data/entities"
	i "server/interfaces"
	"text/template"

	ws "github.com/gorilla/websocket"
)

type ChatRenderer struct {
	authS i.AuthService
	chatS i.ChatService
	msgS  i.MessageService
}

func (cr *ChatRenderer) RenderChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("no session cookie found")
		http.Error(w, "no session cookie found", http.StatusInternalServerError)
		return
	}

	session, err := cr.authS.ValidateAndRefreshSession(cookie.Value)
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

	messages, err := cr.msgS.GetMessages()
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

	testData := []byte("testdata") // testData()

	err = conn.WriteMessage(ws.TextMessage, testData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
		return
	}

	for {
		msgT, _, err := conn.ReadMessage()
		if err != nil {
			break
		} else {
			testEcho := []byte("echo") // testEcho()
			err = conn.WriteMessage(ws.TextMessage, testEcho)
			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
				break
			}
		}
		log.Println(msgT)
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
