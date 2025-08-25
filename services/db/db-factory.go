package dbservice

import (
	i "server/interfaces"
	"server/services/db/SQL"
	prov "server/services/db/providers"
	xerr "server/xerrors"
)

func dbServiceFactory(c prov.Credentials) (i.DbService, error) {
	switch c.Provider() {
	case "sqlite3":
		return SQL.NewDbService(c)
	}
	return nil, xerr.InitDbServiceFail
}
