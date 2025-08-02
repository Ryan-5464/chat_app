package main

import (
	"log"
	"net/http"
	"os"
	"server/handler"
	lgr "server/logging"
	mw "server/middleware"
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

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	dbPath := cwd + "/data/database/sqldb.db"
	log.Println(dbPath)

	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", dbPath)
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
	indexHandler := handler.NewIndexHandler(logger, authS)
	registerHandler := handler.NewRegisterHandler(logger, authS, userS)
	loginHandler := handler.NewLoginHandler(logger, authS, userS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	authMW := mw.NewAuthMiddleware(logger, authS)

	http.Handle("/api/register", authMW.AttachTo(http.HandlerFunc(registerHandler.RegisterUser)))
	http.Handle("/api/login", authMW.AttachTo(http.HandlerFunc(loginHandler.LoginUser)))
	http.Handle("/api/chat/new", authMW.AttachTo(http.HandlerFunc(chatHandler.NewChat)))
	http.Handle("/api/chat/switch", authMW.AttachTo(http.HandlerFunc(chatHandler.SwitchChat)))
	http.Handle("/api/chat/friend/add", authMW.AttachTo(http.HandlerFunc(chatHandler.AddFriend)))
	http.Handle("/login", authMW.AttachTo(http.HandlerFunc(loginHandler.RenderLoginPage)))
	http.Handle("/register", authMW.AttachTo(http.HandlerFunc(registerHandler.RenderRegisterPage)))
	http.Handle("/chat", authMW.AttachTo(http.HandlerFunc(chatHandler.RenderChatPage)))
	http.Handle("/ws", authMW.AttachTo(http.HandlerFunc(chatHandler.ChatWebsocket)))
	http.Handle("/", authMW.AttachTo(http.HandlerFunc(indexHandler.RenderIndexPage)))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
