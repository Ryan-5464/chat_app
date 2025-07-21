package renderers

import (
	"fmt"
	"log"
	"net/http"
	"server/data/entities"
	"text/template"
	"time"

	ws "github.com/gorilla/websocket"
)

type ChatRenderer struct {
}

func (cr *ChatRenderer) RenderChat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/templates/chat.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Chats:    testChats(),
		Messages: testMessages(),
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

func testChats() []entities.Chat {
	chat1 := entities.Chat{
		Id:                 1,
		Name:               "test1",
		AdminId:            3,
		AdminName:          "alf",
		MemberCount:        4,
		UnreadMessageCount: 14,
		CreatedAt:          time.Now(),
	}
	chat2 := entities.Chat{
		Id:                 2,
		Name:               "test2",
		AdminId:            4,
		AdminName:          "derek",
		MemberCount:        3,
		UnreadMessageCount: 2,
		CreatedAt:          time.Now(),
	}
	return []entities.Chat{chat1, chat2}
}

func testMessages() []entities.Message {
	message1 := entities.Message{
		Id:         1,
		UserId:     3,
		ChatId:     1,
		ReplyId:    0,
		Author:     "alf",
		Text:       "hello",
		CreatedAt:  time.Now(),
		LastEditAt: time.Now(),
	}
	message2 := entities.Message{
		Id:         2,
		UserId:     3,
		ChatId:     1,
		ReplyId:    0,
		Author:     "alf",
		Text:       "there",
		CreatedAt:  time.Now(),
		LastEditAt: time.Now(),
	}
	return []entities.Message{message1, message2}
}
