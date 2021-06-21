package db

type Type string

const(
	POSTGRES Type = "postgresql"
	MONGO    Type = "mongodb"
)

type Properties struct {
	Type     Type   `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"database"`
	Host     string `json:"host"`
	Port     uint16  `json:"port"`
}

type Option func(db *Properties)


/// Methods

func WithUser(username string) Option {
	return func(db *Properties) {
		db.Username = username
	}
}

func WithPass(password string) Option {
	return func(db *Properties) {
		db.Password = password
	}
}

func WithDBName(dbname string) Option {
	return func(db *Properties) {
		db.Dbname = dbname
	}
}

func WithHost(host string) Option {
	return func(db *Properties) {
		db.Host = host
	}
}

func WithPort(port uint16) Option {
	return func(db *Properties) {
		db.Port = port
	}
}
