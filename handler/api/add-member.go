package api

import (
	"encoding/json"
	"net/http"
	ent "server/data/entities"
	mw "server/handler/middleware"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	typ "server/types"
	"server/util"
)

func AddMember(a i.AuthService, c i.ChatService, cn i.ConnectionService, m i.MessageService) http.Handler {
	h := addMember{
		chatS: c,
		connS: cn,
		msgS:  m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type addMember struct {
	chatS i.ChatService
	connS i.ConnectionService
	msgS  i.MessageService
}

func (h addMember) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	var req amrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	util.Log.Dbugf("Add member email: %v", req.Email)
	util.Log.Dbugf("Chat id: %v", req.ChatId)

	res, err := h.handleRequest(req)
	if err != nil {
		util.Log.Errorf("failed to handle add member request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")

}

func (h addMember) handleRequest(req amrequest) (amresponse, error) {
	util.Log.FunctionInfo()

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return amresponse{}, err
	}

	memberId, err := h.chatS.AddMember(cred.Email(req.Email), chatId)
	if err != nil {
		return amresponse{}, err
	}

	util.Log.Dbugf("NEW MEMBER ID => %v", memberId)

	member, err := h.chatS.GetChatMember(chatId, memberId)
	if err != nil {
		return amresponse{}, err
	}

	lastestMsgId, err := h.msgS.GetLatestMessageId()
	if err != nil {
		return amresponse{}, err
	}

	if err := h.msgS.UpdateLastReadMsgId(lastestMsgId, chatId, memberId); err != nil {
		return amresponse{}, err
	}

	chats, err := h.chatS.GetChats(memberId)
	if err != nil {
		return amresponse{}, err
	}

	conn := h.connS.GetConnection(memberId)

	payload := struct {
		Type  string
		Chats []ent.Chat
	}{
		Type:  "AddMember",
		Chats: chats,
	}

	if err := conn.WriteJSON(payload); err != nil {
		util.Log.Errorf("failed to write to websocket connection: %v", err)
		return amresponse{}, err
	}

	res := amresponse{
		Members: []ent.Member{*member},
	}

	return res, nil
}

type amrequest struct {
	Email  string `json:"Email"`
	ChatId string `json:"ChatId"`
}

type amresponse struct {
	Members []ent.Member
}
