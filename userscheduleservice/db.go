package userscheduleservice

import (
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DSN string
}

type SqlDb struct {
	*sqlx.DB
}

func ProvideDB(cfg *Config) SqlDb {
	dsn := cfg.DSN
	var err error
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return SqlDb{db}
}
