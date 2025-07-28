package providers

type DbProvider string

func (d DbProvider) String() string {
	return string(d)
}

type DbPath string

func (d DbPath) String() string {
	return string(d)
}

const (
	SQLite3    DbProvider = "sqlite3"
	PostgreSQL DbProvider = "postgreSQL"
	InMemoryDb DbPath     = ":memory:"
)

type DbConfig map[string]string

func (d DbConfig) Add(key string, value string) {
	(d)[key] = value
}

func (d DbConfig) Get(key string) string {
	return d[key]
}

func NewDbCredentials(dbProvider DbProvider, config DbConfig) Credentials {
	return Credentials{
		dbProvider: dbProvider,
		config:     config,
	}
}

type Credentials struct {
	dbProvider DbProvider
	config     DbConfig
}

func (c Credentials) Provider() string {
	return string(c.dbProvider)
}

func (c Credentials) Value(key string) string {
	return c.config[key]
}
