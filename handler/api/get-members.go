package api

import (
	"net/http"
	ent "server/data/entities"
	mw "server/handler/middleware"
	i "server/interfaces"
	typ "server/types"
	"server/util"
)

func GetMembers(a i.AuthService, c i.ChatService) http.Handler {
	h := getMembers{
		chatS: c,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type getMembers struct {
	chatS i.ChatService
}

func (h getMembers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	query := r.URL.Query()
	req := gmrequest{
		ChatId: query.Get("ChatId"),
	}

	res, err := h.handleRequest(req)
	if err != nil {
		util.Log.Errorf("failed to handle switch chat gmrequest, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h getMembers) handleRequest(req gmrequest) (gmresponse, error) {
	util.Log.FunctionInfo()

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return gmresponse{}, err
	}

	members, err := h.chatS.GetChatMembers(chatId)
	if err != nil {
		return gmresponse{}, err
	}

	res := gmresponse{
		Members: members,
	}

	return res, nil
}

type gmrequest struct {
	ChatId string `json:"ChatId"`
}

type gmresponse struct {
	Members []ent.Member
}
