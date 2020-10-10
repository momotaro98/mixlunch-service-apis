package userservice

import (
	"database/sql"
)

type Config struct {
	DSN string
}

type SqlDb struct {
	*sql.DB
}

func ProvideDB(cfg *Config) SqlDb {
	dsn := cfg.DSN
	var err error
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return SqlDb{db}
}
