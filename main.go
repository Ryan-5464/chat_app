package main

import (
	"log"
	"net/http"
	"os"
	"server/handler/api"
	"server/handler/socket"
	"server/handler/view"
	sauth "server/services/auth"
	schat "server/services/chat"
	sconn "server/services/conn"
	dbs "server/services/db"
	prov "server/services/db/providers"
	smsg "server/services/message"
	repo "server/services/repository"
	suser "server/services/user"
)

func main() {

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

	dbService, err := dbs.NewDbService(c)
	if err != nil {
		log.Println(err)
		return
	}

	userR := repo.NewUserRepository(dbService)
	chatR := repo.NewChatRepository(dbService)
	msgR := repo.NewMessageRepository(dbService)

	authS := sauth.NewAuthService()
	connS := sconn.NewConnectionService()
	userS := suser.NewUserService(userR)
	chatS := schat.NewChatService(chatR, userS)
	msgS := smsg.NewMessageService(msgR, userS, connS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/api/chat/members", api.GetMembers(authS, chatS))
	http.Handle("/api/chat/switch", api.SwitchChat(authS, msgS))
	http.Handle("/api/chat/new", api.NewChat(authS, chatS))
	http.Handle("/api/chat/member/remove", api.RemoveMember(authS, chatS))
	http.Handle("/api/chat/contact/remove", api.RemoveContact(authS, chatS, msgS, userS))
	http.Handle("/api/chat/leave", api.LeaveChat(authS, chatS, msgS))
	http.Handle("/api/chat/edit", api.EditChatName(authS, chatS))
	http.Handle("/api/chat/contact/switch", api.SwitchContactChat(authS, msgS))
	http.Handle("/api/chat/contact/add", api.AddContact(authS, connS, userS))
	http.Handle("/api/chat/members/add", api.AddMember(authS, chatS))

	http.Handle("/api/message/delete", api.DeleteMessage(authS, msgS))
	http.Handle("/api/message/edit", api.EditMessage(authS, msgS))

	http.Handle("/api/profile/name/edit", api.EditUserName(authS, userS))
	http.Handle("/api/register", api.Register(authS, userS))
	http.Handle("/api/login", api.Login(authS, userS))

	http.Handle("/profile", view.Profile(authS, userS))
	http.Handle("/chat", view.Chat(authS, chatS, msgS, userS))
	http.Handle("/ws", socket.Chat(authS, connS, msgS))
	http.Handle("/register", view.Register())
	http.Handle("/login", view.Login())
	http.Handle("/", view.Index())

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
