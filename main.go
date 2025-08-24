package main

import (
	"log"
	"net/http"
	"os"
	"server/handler"
	"server/handler/api"
	"server/handler/socket"
	"server/handler/view"
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
	logger := lgr.NewLogger(false)

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	dbPath := cwd + "/data/database/app.db"
	log.Println(dbPath)

	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", dbPath)
	c := prov.NewDbCredentials(prov.SQLite3, config)

	dbService, err := dbs.NewDbService(logger, c)
	if err != nil {
		log.Println(err)
		return
	}

	userR := repo.NewUserRepository(logger, dbService)
	chatR := repo.NewChatRepository(logger, dbService)
	msgR := repo.NewMessageRepository(logger, dbService)

	authS := sauth.NewAuthService(logger)
	connS := sconn.NewConnectionService(logger)
	userS := suser.NewUserService(logger, userR)
	chatS := schat.NewChatService(logger, chatR, userS)
	msgS := smsg.NewMessageService(logger, msgR, userS, connS)

	chatHandler := handler.NewChatHandler(logger, authS, chatS, msgS, connS, userS)
	indexHandler := handler.NewIndexHandler(logger, authS)
	registerHandler := handler.NewRegisterHandler(logger, authS, userS)
	loginHandler := handler.NewLoginHandler(logger, authS, userS)
	profileHandler := handler.NewProfileHandler(logger, authS, userS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/api/chat/members", api.GetMembers(authS, chatS))
	http.Handle("/api/chat/switch", api.SwitchChat(authS, msgS))
	http.Handle("/api/chat/new", api.NewChat(authS, chatS))
	http.Handle("/api/chat/member/remove", api.RemoveMember(authS, chatS))
	http.Handle("/api/chat/contact/remove", api.RemoveContact(authS, chatS, msgS))
	http.Handle("/api/chat/leave", api.LeaveChat(authS, chatS, msgS))

	http.Handle("/api/message/delete", api.DeleteMessage(authS, msgS))

	http.Handle("/chat", view.Chat(authS, chatS, msgS, userS))

	http.Handle("/ws", socket.Chat(authS, chatS, msgS))

	http.Handle("/api/register", authMW.AttachTo(http.HandlerFunc(registerHandler.RegisterUser)))
	http.Handle("/api/login", authMW.AttachTo(http.HandlerFunc(loginHandler.LoginUser)))
	http.Handle("/api/chat/edit", authMW.AttachTo(http.HandlerFunc(chatHandler.EditChatName)))
	http.Handle("/api/chat/contact/switch", authMW.AttachTo(http.HandlerFunc(chatHandler.SwitchContactChat)))
	http.Handle("/api/chat/contact/add", authMW.AttachTo(http.HandlerFunc(chatHandler.AddContact)))
	http.Handle("/api/chat/members/add", authMW.AttachTo(http.HandlerFunc(chatHandler.AddMemberToChat)))
	http.Handle("/api/message/edit", authMW.AttachTo(http.HandlerFunc(chatHandler.EditMessage)))
	http.Handle("/api/profile/name/edit", authMW.AttachTo(http.HandlerFunc(profileHandler.EditUserName)))

	http.Handle("/profile", authMW.AttachTo(http.HandlerFunc(profileHandler.RenderProfilePage)))
	http.Handle("/login", authMW.AttachTo(http.HandlerFunc(loginHandler.RenderLoginPage)))
	http.Handle("/register", authMW.AttachTo(http.HandlerFunc(registerHandler.RenderRegisterPage)))
	http.Handle("/chat", authMW.AttachTo(reqMethodMW.AttachTo(http.HandlerFunc(chatHandler.RenderChatPage), mw.GET)))
	http.Handle("/", authMW.AttachTo(http.HandlerFunc(indexHandler.RenderIndexPage)))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
