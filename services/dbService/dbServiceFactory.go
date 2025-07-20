package dbservice

import (
	i "server/interfaces"
	"server/services/dbService/SQL"
	typ "server/types"
	xerr "server/xerrors"
)

func dbServiceFactory(c typ.Credentials) (i.DbService, error) {
	switch c.Provider() {
	case "sqlite3":
		return SQL.NewDbService(c)
	}
	return nil, xerr.InitDbServiceFail
}
