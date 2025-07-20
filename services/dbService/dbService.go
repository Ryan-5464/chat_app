package dbservice

import (
	"fmt"
	i "server/interfaces"
	typ "server/types"
)

func NewDbService(c typ.Credentials) (i.DbService, error) {
	dbS, err := dbServiceFactory(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dbService: %w", err)
	}
	return dbS, nil
}
