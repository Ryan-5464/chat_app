package main

import (
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
	chatRepo := repo.NewChatRepository(dbService)
	messageRepo := repo.NewMessageRepository(dbService)
	userRepo.GetUsers()
}
