package templ

import (
	"fmt"
	"text/template"
)

type HTMLPath = string

const (
	chatViewHTML        HTMLPath = "./static/templates/chat/chat-view.html"
	messagesHTML        HTMLPath = "./static/templates/chat/messages.html"
	messageHTML         HTMLPath = "./static/templates/chat/message.html"
	chatHTML            HTMLPath = "./static/templates/chat/chat.html"
	chatsHTML           HTMLPath = "./static/templates/chat/chats.html"
	newChatHTML         HTMLPath = "./static/templates/chat/new-chat.html"
	newChatNameHTML     HTMLPath = "./static/templates/chat/new-chat-name.html"
	LeaveChatHTML       HTMLPath = "./static/templates/chat/leave-chat.html"
	contactHTML         HTMLPath = "./static/templates/chat/contact.html"
	contactsHTML        HTMLPath = "./static/templates/chat/contacts.html"
	contactModalHTML    HTMLPath = "./static/templates/chat/contact-modal.html"
	chatModalHTML       HTMLPath = "./static/templates/chat/chat-modal.html"
	messageModalHTML    HTMLPath = "./static/templates/chat/message-modal.html"
	memberListModalHTML HTMLPath = "./static/templates/chat/member-list-modal.html"
	memberModalHTML     HTMLPath = "./static/templates/chat/member-modal.html"
)

var (
	ChatView *template.Template
)

func init() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": dict,
	})
	ChatView = template.Must(
		tmpl.ParseFiles(
			chatViewHTML,
			messagesHTML,
			messageHTML,
			chatHTML,
			chatsHTML,
			newChatHTML,
			newChatNameHTML,
			LeaveChatHTML,
			contactHTML,
			contactsHTML,
			contactModalHTML,
			chatModalHTML,
			messageModalHTML,
			memberListModalHTML,
			memberModalHTML,
		),
	)
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: must have even number of args")
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}
