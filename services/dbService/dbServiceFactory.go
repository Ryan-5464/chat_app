package dbservice

import (
	i "server/interfaces"
	"server/services/dbService/SQL"
	prov "server/services/dbService/providers"
	xerr "server/xerrors"
)

func dbServiceFactory(lgr i.Logger, c prov.Credentials) (i.DbService, error) {
	switch c.Provider() {
	case "sqlite3":
		return SQL.NewDbService(lgr, c)
	}
	return nil, xerr.InitDbServiceFail
}
