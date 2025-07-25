package dbservice

import (
	"fmt"
	i "server/interfaces"
	prov "server/services/dbService/providers"
)

func NewDbService(lgr i.Logger, c prov.Credentials) (i.DbService, error) {
	dbS, err := dbServiceFactory(lgr, c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dbService: %w", err)
	}
	return dbS, nil
}
