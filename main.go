package main

import (
	"log"
	"net/http"
	"server/handlers"
	"server/handlers/renderers"
	lgr "server/logging"
	sauth "server/services/authService"
	schat "server/services/chatService"
	sconn "server/services/connService"
	dbs "server/services/dbService"
	prov "server/services/dbService/providers"
	smsg "server/services/messageService"
	repo "server/services/repository"
	suser "server/services/userService"
)

func main() {
	logger := lgr.NewLogger(true)

	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", ":memory:")
	c := prov.NewCredentials(prov.SQLite3, config)

	dbService, err := dbs.NewDbService(logger, c)
	if err != nil {
		return
	}

	userR := repo.NewUserRepository(logger, dbService)
	chatR := repo.NewChatRepository(logger, dbService)
	msgR := repo.NewMessageRepository(logger, dbService)

	authS := sauth.NewAuthService(logger)
	chatS := schat.NewChatService(logger, chatR)
	connS := sconn.NewConnectionService(logger)
	userS := suser.NewUserService(logger, userR)
	msgS := smsg.NewMessageService(logger, msgR, userS, connS)

	cr := renderers.NewChatRenderer(logger, authS, chatS, msgS, connS, userS)
	testRegistrationHandler := handlers.NewTestRegistrationHandler(logger, userS)
	testChatHandler := handlers.NewTestChatHandler(logger, chatS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/r", testRegistrationHandler.RegisterUser)
	http.HandleFunc("/chat", cr.RenderChat)
	http.HandleFunc("/x", testChatHandler.NewChat)
	http.HandleFunc("/ws", cr.ChatWebsocket)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
