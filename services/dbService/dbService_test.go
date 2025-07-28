package dbservice

import (
	i "server/interfaces"
	"server/services/dbService/SQL"
	prov "server/services/dbService/providers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testConfig() prov.DbConfig {
	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", ":memory:")
	return config
}

func TestNewCredentials(t *testing.T) {
	c := prov.NewDbCredentials(prov.SQLite3, testConfig())
	assert.Equal(t, c.Provider(), prov.SQLite3.String())
	assert.Equal(t, c.Value("driver"), "sqlite3")
	assert.Equal(t, c.Value("path"), ":memory:")
}

func TestDbServiceFactory(t *testing.T) {
	c := prov.NewDbCredentials(prov.SQLite3, testConfig())
	dbS, err := dbServiceFactory(nil, c)
	if err != nil {
		t.Errorf("failed to initialize dbService %v", err)
	}

	var _ i.DbService = (*SQL.DbService)(nil)
	if _, ok := any(dbS).(i.DbService); !ok {
		t.Fatal("DatabaseService does not implement DbService interface")
	}
}
