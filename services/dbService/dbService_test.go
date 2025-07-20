package dbservice

import (
	"server/data/entities"
	typ "server/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDetails() map[string]string {
	details := make(map[string]string)
	details["driver"] = "sqlite3"
	details["path"] = ":memory:"
	return details
}

func TestNewCredentials(t *testing.T) {
	c := typ.NewCredentials(typ.SQLite3, testDetails())
	assert.Equal(t, c.Provider(), string(typ.SQLite3))
	assert.Equal(t, c.Value("driver"), "sqlite3")
	assert.Equal(t, c.Value("path"), ":memory:")
}

func TestDbServiceFactory(t *testing.T) {
	c := typ.NewCredentials(typ.SQLite3, testDetails())
	dbS, err := dbServiceFactory(c)
	if err != nil {
		t.Errorf("failed to initialize dbService %v", err)
	}
	users := dbS.GetUsers()
	assert.Equal(t, []entities.User{}, users)
}
