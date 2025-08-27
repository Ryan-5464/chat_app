package view

import (
	"net/http"
	ent "server/data/entities"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	"server/handler/templ"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func Chat(a i.AuthService, c i.ChatService, m i.MessageService, u i.UserService, cn i.ConnectionService) http.Handler {
	h := chatView{
		chatS: c,
		msgS:  m,
		userS: u,
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type chatView struct {
	chatS i.ChatService
	msgS  i.MessageService
	userS i.UserService
	connS i.ConnectionService
}

func (h chatView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	d, err := h.getChatTemplateData(session.UserId())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := templ.ChatView.ExecuteTemplate(w, "chatView", d); err != nil {
		util.Log.Errorf("failed to execute chatView template, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h chatView) getChatTemplateData(userId typ.UserId) (tmplData, error) {
	util.Log.FunctionInfo()

	chats, err := h.chatS.GetChats(userId)
	if err != nil {
		util.Log.Errorf("failed to get chats for userId: %v, Error: %v", userId, err)
		return tmplData{}, err
	}

	var chatId typ.ChatId
	var messages []ent.Message
	if len(chats) != 0 {
		chatId = chats[0].Id
		messages, err = h.msgS.GetChatMessages(chatId, userId)
		if err != nil {
			util.Log.Errorf("failed to get chat messages for chatId: %v, Error: %v", chatId, err)
			return tmplData{}, err
		}
	}
	// Setting to zero since active chat will always display latest message
	// Must do this because chats are retrieved before latest message id is updated so active chatid
	// is outdated.

	if len(chats) != 0 {
		for i := range chats {
			chats[i].UnreadMessageCount, err = h.chatS.GetUnreadMessageCount(chats[i].Id, userId)
			if err != nil {
				return tmplData{}, err
			}
		}
		chats[0].UnreadMessageCount = 0
	}

	latestMsgId := findLastestMessageId(messages)

	if err := h.msgS.UpdateLastReadMsgId(latestMsgId, chatId, userId); err != nil {
		return tmplData{}, err
	}

	contacts, err := h.userS.GetContacts(userId)
	if err != nil {
		util.Log.Errorf("failed to get contacts for userId: %v, Error: %v", userId, err)
		return tmplData{}, err
	}

	user, err := h.userS.GetUser(userId)
	if err != nil {
		util.Log.Errorf("failed to get contacts for userId: %v, Error: %v", userId, err)
		return tmplData{}, err
	}

	return tmplData{
		User:         user,
		UserId:       userId,
		Chats:        chats,
		Messages:     messages,
		Contacts:     contacts,
		ActiveChatId: chatId,
	}, nil

}

type tmplData struct {
	User         *ent.User     `json:"User"`
	UserId       typ.UserId    `json:"UserId"`
	Chats        []ent.Chat    `json:"Chats"`
	Messages     []ent.Message `json:"Messages"`
	Contacts     []ent.Contact `json:"Contacts"`
	ActiveChatId typ.ChatId    `json:"ActiveChatId"`
}
