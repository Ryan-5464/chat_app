package socket

import (
	"net/http"
	"server/util"

	"github.com/gorilla/websocket"
)

func newSocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Log.Errorf("Failed to establish websocket connection! => error: %v", err)
		return nil
	}

	return conn
}
