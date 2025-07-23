package main

import (
	"log"
	"net/http"
	"server/handlers/renderers"
	sauth "server/services/authService"
	schat "server/services/chatService"
	dbs "server/services/dbService"
	prov "server/services/dbService/providers"
	smsg "server/services/messageService"
	repo "server/services/repository"
)

func main() {

	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", ":memory:")
	c := prov.NewCredentials(prov.SQLite3, config)

	dbService, err := dbs.NewDbService(c)
	if err != nil {
		return
	}

	userR := repo.NewUserRepository(dbService)
	chatR := repo.NewChatRepository(dbService)
	msgR := repo.NewMessageRepository(dbService)
	userR.GetUsers()

	authS := sauth.NewAuthService()
	chatS := schat.NewChatService(chatR)
	msgS := smsg.NewMessageService(msgR)

	cr := renderers.NewChatRenderer(authS, chatS, msgS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", cr.RenderChat)
	http.HandleFunc("/ws", cr.ChatWebsocket)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
