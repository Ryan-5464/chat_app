package socket

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
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

func Chat(a i.AuthService, c i.ConnectionService, m i.MessageService) http.Handler {
	h := chatWebSocket{
		connS: c,
		msgS:  m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type chatWebSocket struct {
	connS i.ConnectionService
	msgS  i.MessageService
}

func (h chatWebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := newSocket(w, r)
	if conn == nil {
		util.Log.Error(errors.New("failed to connect user"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)
	userId := session.UserId()

	h.connS.StoreConnection(conn, userId)
	defer h.connS.DisconnectUser(userId)

	util.Log.Dbugf("active connections : %v", h.connS.GetActiveConnections())

	for {

		payload, err := h.readIncomingMessage(conn)
		if err != nil {
			util.Log.Errorf("failed to read incoming websocket message, Error: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			break
		}

		switch payload.Type {
		case NewMessage:
			util.Log.Dbug("Handling new message...")

			req, err := payload.ParseNewMessage()
			if err != nil {
				util.Log.Errorf("Failed to parse message request %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				break
			}

			if err = h.handleNewMessageRequest(req, userId); err != nil {
				util.Log.Errorf("Failed to handle message request %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				break
			}

		case NewContactMessage:
			util.Log.Dbug("Handling new contact message...")

			req, err := payload.ParseNewMessage()
			if err != nil {
				util.Log.Errorf("Failed to parse message request %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				break
			}

			if err = h.handleNewContactMessageRequest(req, userId); err != nil {
				util.Log.Errorf("Failed to handle message request %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				break
			}

		}
		util.Log.Dbug("->>>> RESPONSE SENT")

	}

	util.Log.Dbug("User disconnected")
}

func (h chatWebSocket) readIncomingMessage(conn i.Socket) (websocketPayload, error) {
	util.Log.FunctionInfo()

	payload := websocketPayload{}
	err := conn.ReadJSON(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (h chatWebSocket) handleNewMessageRequest(req request, userId typ.UserId) error {
	util.Log.FunctionInfo()

	util.Log.Dbugf("chatId %v: replyId %v", req.ChatId, req.ReplyId)
	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return err
	}

	var replyId typ.MessageId
	if req.ReplyId != "" {
		replyId, err = typ.ToMessageId(req.ReplyId)
		if err != nil {
			return err
		}
	}

	err = h.msgS.HandleNewMessage(userId, chatId, &replyId, req.MsgText)
	if err != nil {
		return err
	}

	return nil
}

func (h chatWebSocket) handleNewContactMessageRequest(req request, userId typ.UserId) error {
	util.Log.FunctionInfo()

	util.Log.Dbugf("chatId %v: replyId %v", req.ChatId, req.ReplyId)
	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return err
	}

	var replyId typ.MessageId
	if req.ReplyId != "" {
		replyId, err = typ.ToMessageId(req.ReplyId)
		if err != nil {
			return err
		}
	}

	err = h.msgS.HandleNewContactMessage(userId, chatId, &replyId, req.MsgText)
	if err != nil {
		return err
	}

	return nil
}

type websocketPayload struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}

func (w *websocketPayload) ParseNewMessage() (request, error) {
	req := request{}
	if err := json.Unmarshal(w.Data, &req); err != nil {
		return request{}, err
	}
	return req, nil
}

type request struct {
	ChatId  string `json:"ChatId"`
	ReplyId string `json:"ReplyId"`
	MsgText string `json:"MsgText"`
}
