package database

import (
	"database/sql"
	"errors"
	"fmt"
	"server/logging"
	"server/services/dbService/SQL/schema"
	"server/services/dbService/providers"
	typ "server/types"
	"testing"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	errUserCreationFail string = "failed to create new user:"
	errReadUserFail     string = "failed to read user:"
	errUpdateUserFail   string = "failed to update user:"
	errDeleteUserFail   string = "failed to delete user:"
)

func TestDBCreate_Success(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	_, err := createUser(db)
	if err != nil {
		t.Fatalf("%s %v", errUserCreationFail, err)
	}
}

func TestDBCreate_ErrConstraintUnique(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	for i := 0; i < 3; i++ {

		_, err := createUser(db)
		if err != nil {
			if i == 1 {
				var sqliteErr sqlite3.Error
				if errors.As(err, &sqliteErr) {
					require.Equal(t, sqlite3.ErrConstraintUnique, sqliteErr.ExtendedCode)
				} else {
					t.Fatalf("expected sqlite3.Error, got: %T: %v", err, err)
				}
				return
			}
			t.Fatalf("%s %v", errUserCreationFail, err)
		}
	}
}

func TestRead_Success(t *testing.T) {
	db, res := initializeTest(t)
	defer db.Close()

	userId, _ := res.LastInsertId()
	rows, err := readUser(t, db, userId)
	if err != nil {
		t.Fatalf("%s %v", errReadUserFail, err)
	}

	for _, row := range rows {
		userName := row[schema.Name].(string)
		email := row[schema.Email].(string)
		pwdHash := row[schema.PwdHash].(string)
		assert.Equal(t, "testuser", userName)
		assert.Equal(t, "testemail@outlook.com", email)
		assert.Equal(t, "testmessage", pwdHash)
	}
}

func TestRead_FailNoRowFound(t *testing.T) {
	db, _ := initializeTest(t)
	defer db.Close()

	userIdNotexist := int64(2000)
	rows, err := readUser(t, db, userIdNotexist)
	if err != nil {
		t.Fatalf("%s %v", errReadUserFail, err)
	}

	var expected typ.Rows
	assert.Equal(t, expected, rows)
}

func TestUpdate_Success(t *testing.T) {
	db, res := initializeTest(t)
	defer db.Close()

	query := fmt.Sprintf(
		"UPDATE %s SET %s = ? WHERE %s = ?",
		schema.UserTable, schema.Email, schema.UserId,
	)

	email := "updatedemail@outlook.com"
	userId, _ := res.LastInsertId()
	if err := db.Update(query, email, userId); err != nil {
		t.Fatalf("%s %v", errUpdateUserFail, err)
	}

	rows, err := readUser(t, db, userId)
	if err != nil {
		t.Fatalf("%s %v", errReadUserFail, err)
	}

	var updatedEmail string
	for _, row := range rows {
		updatedEmail = row[schema.Email].(string)
	}

	assert.Equal(t, "updatedemail@outlook.com", updatedEmail)
}

func TestDelete_Success(t *testing.T) {
	db, res := initializeTest(t)
	defer db.Close()

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?",
		schema.UserTable, schema.UserId,
	)

	userId, _ := res.LastInsertId()
	if err := db.Delete(query, userId); err != nil {
		t.Fatalf("%s %v", errDeleteUserFail, err)
	}

	rows, err := readUser(t, db, userId)
	if err != nil {
		require.Error(t, err)
		return
	}
	var expected typ.Rows
	assert.Equal(t, expected, rows)
}

func TestDelete_FailRowNotExist(t *testing.T) {
	db, _ := initializeTest(t)
	defer db.Close()

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = ?",
		schema.UserTable, schema.UserId,
	)

	userIdNotExist := 2000
	if err := db.Delete(query, userIdNotExist); err != nil {
		t.Fatalf("%s %v", errDeleteUserFail, err)
	}

	rows, err := readUser(t, db, 1)
	if err != nil {
		require.Error(t, err)
		return
	}
	var expected typ.Rows
	assert.NotEqual(t, expected, rows)
}

func initializeTest(t *testing.T) (*DB, sql.Result) {
	db := openDB(t)

	res, err := createUser(db)
	if err != nil {
		t.Fatalf("%s %v", errUserCreationFail, err)
	}

	return db, res
}

func openDB(t *testing.T) *DB {
	config := providers.DbConfig{}
	config.Add("driver", providers.SQLite3.String())
	config.Add("path", providers.InMemoryDb.String())
	dbCredentials := providers.NewDbCredentials(providers.SQLite3, config)
	logger := logging.NewLogger(false)
	db, err := NewDatabase(logger, dbCredentials)
	if err != nil {
		t.Fatalf("Database creation failed: %v", err)
	}
	return db
}

func createUser(db *DB) (sql.Result, error) {
	userName := "testuser"
	email := "testemail@outlook.com"
	pwdHash := "testmessage"

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)",
		schema.UserTable,
		schema.Name,
		schema.Email,
		schema.PwdHash,
	)

	return db.Create(query, userName, email, pwdHash)
}

func readUser(t *testing.T, db *DB, userId int64) (typ.Rows, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = ?",
		schema.UserTable, schema.UserId,
	)

	rows, err := db.Read(query, userId)
	if err != nil {
		t.Fatalf("%s %v", errReadUserFail, err)
	}

	return rows, nil

}
