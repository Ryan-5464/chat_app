package main

import (
	"log"
	"net/http"
	"server/handler"
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
	c := prov.NewDbCredentials(prov.SQLite3, config)

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

	chatHandler := handler.NewChatHandler(logger, authS, chatS, msgS, connS, userS)
	indexHandler := handler.NewIndexHandler(logger)
	registerHandler := handler.NewRegisterHandler(logger, authS, userS)
	loginHandler := handler.NewLoginHandler(logger, authS, userS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/register", registerHandler.RenderRegisterPage)
	http.HandleFunc("/api/register", registerHandler.RegisterUser)
	http.HandleFunc("/login", loginHandler.RenderLoginPage)
	http.HandleFunc("/api/login", loginHandler.LoginUser)
	http.HandleFunc("/chat", chatHandler.RenderChatPage)
	http.HandleFunc("/ws", chatHandler.ChatWebsocket)
	http.HandleFunc("/", indexHandler.RenderIndexPage)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
