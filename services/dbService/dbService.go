package dbservice

import (
	"fmt"
	i "server/interfaces"
	prov "server/services/dbService/providers"
)

func NewDbService(c prov.Credentials) (i.DbService, error) {
	dbS, err := dbServiceFactory(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dbService: %w", err)
	}
	return dbS, nil
}
