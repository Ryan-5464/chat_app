package handler

import (
	"log"
	"net/http"
	td "server/data/test"
	i "server/interfaces"
)

func NewTestChatHandler(l i.Logger, c i.ChatService) *TestChatHandler {
	return &TestChatHandler{
		lgr:   l,
		chatS: c,
	}
}

type TestChatHandler struct {
	lgr   i.Logger
	chatS i.ChatService
}

func (t *TestChatHandler) NewChat(w http.ResponseWriter, r *http.Request) {
	t.lgr.LogFunctionInfo()

	chat := td.TestChat()
	newChat, err := t.chatS.NewChat(chat)
	if err != nil {
		log.Println("chat creation failed", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("new chat: ", newChat)

}
