package renderers

import (
	"fmt"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type ChatRenderer struct {
}

func (cr *ChatRenderer) RenderChat(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/templates/chat.html")
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

// func testEcho() []byte {

// 	message1 := testMessage{
// 		MessageId:   1,
// 		MessageText: "echo",
// 	}
// 	messages := []testMessage{message1}

// 	chat1 := testChat{
// 		ChatId:   1,
// 		ChatName: "first",
// 	}
// 	chat2 := testChat{
// 		ChatId:   2,
// 		ChatName: "second",
// 	}
// 	chats := []testChat{chat1, chat2}

// 	payload := testPayload{
// 		Messages: messages,
// 		Chats:    chats,
// 	}

// 	data, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil
// 	}

// 	return data
// }
