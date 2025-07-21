package main

import (
	"log"
	"net/http"
	"server/handlers/renderers"
	dbs "server/services/dbService"
	prov "server/services/dbService/providers"
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

	userRepo := repo.NewUserRepository(dbService)
	// chatRepo := repo.NewChatRepository(dbService)
	// messageRepo := repo.NewMessageRepository(dbService)
	userRepo.GetUsers()

	cr := renderers.ChatRenderer{}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", cr.RenderChat)
	http.HandleFunc("/ws", cr.ChatWebsocket)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
