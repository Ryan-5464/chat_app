package types

type ChatId int64

func (c ChatId) Int64() int64 {
	return int64(c)
}

type MessageId int64

func (m MessageId) Int64() int64 {
	return int64(m)
}

type UserId int64

func (u UserId) Int64() int64 {
	return int64(u)
}

type Rows []map[string]any

type DbProvider string

func (d DbProvider) String() string {
	return string(d)
}

const (
	SQLite3    DbProvider = "sqlite3"
	PostgreSQL DbProvider = "postgreSQL"
)

func NewCredentials(dbProvider DbProvider, details map[string]string) Credentials {
	return Credentials{
		dbProvider: dbProvider,
		details:    details,
	}
}

type Credentials struct {
	dbProvider DbProvider
	details    map[string]string
}

func (c Credentials) Provider() string {
	return string(c.dbProvider)
}

func (c Credentials) Value(key string) string {
	return c.details[key]
}
