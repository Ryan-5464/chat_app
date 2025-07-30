package socket

import (
	"fmt"
	"net/http"
	l "server/logging"

	"github.com/gorilla/websocket"
)

func New(w http.ResponseWriter, r *http.Request) *websocket.Conn {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errWebsocketOpenFail := fmt.Errorf("Failed to establish websocket connection!")
		l.Lgr.LogError(errWebsocketOpenFail)
		return nil
	}

	return conn
}
