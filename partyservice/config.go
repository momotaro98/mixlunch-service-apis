package partyservice

import (
	"context"
	"database/sql"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Config struct {
	DSN                   string
	AppCredentialFilePath string
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

type App struct {
	*firebase.App
}

func ProvideApp(cnf *Config) App {
	var err error
	opt := option.WithCredentialsFile(cnf.AppCredentialFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	return App{app}
}
