package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/momotaro98/mixlunch-service-api/logger"
	"github.com/momotaro98/mixlunch-service-api/partyservice"
	"github.com/momotaro98/mixlunch-service-api/tagservice"
	usService "github.com/momotaro98/mixlunch-service-api/userscheduleservice"
	"github.com/momotaro98/mixlunch-service-api/userservice"
)

const (
	Version               = "1.3.1"
	ServiceAccountKeyPath = "./serviceAccount/serviceAccountKey.json"
)

func main() {
	var (
		httpAddr = flag.String("http", ":5000", "http listen address")
	)
	flag.Parse()

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// Initial loading config
	var (
		logConf = &logger.Config{
			ErrorLevel: logger.Info, // Might need to be set from command argument
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
		)
		usConf = &usService.Config{
			DSN: dsn,
		}
		pConf = &partyservice.Config{
			DSN:                   dsn,
			AppCredentialFilePath: ServiceAccountKeyPath,
		}
		tConf = &tagservice.Config{
			DSN: dsn,
		}
		uConf = &userservice.Config{
			DSN: dsn,
		}
	)

	// Middlewares

	// Auth middleware
	var authActivate bool
	if os.Getenv("AUTH_ACTIVATE") == "ON" {
		authActivate = true
	}
	auth := AuthMiddle(authActivate, logConf, ServiceAccountKeyPath)

	const (
		GET  = "GET"
		POST = "POST"
	)

	// Launch REST server
	go func() {
		r := mux.NewRouter()
		s := r.PathPrefix("/api/v1").Subrouter()

		// User schedule server
		s.Handle("/userschedule/{uid:[a-zA-Z0-9]+}/{beginDateTime}/{endDateTime}",
			M(initializeUserScheduleHandler(logConf, usConf, tConf), auth)).
			Methods(GET)
		// [Note] The order of the p.Add routing is crucial
		// "update" must be before none("add") because "update" is also a case of {uid:[a-zA-Z0-9]+}
		s.Handle("/userschedule/update/{uid:[a-zA-Z0-9]+}/",
			M(initializeUpdateUserScheduleHandler(logConf, usConf, tConf), auth)).
			Methods(POST)
		s.Handle("/userschedule/delete/{uid:[a-zA-Z0-9]+}/",
			M(initializeDeleteUserScheduleHandler(logConf, usConf, tConf), auth)).
			Methods(POST)
		s.Handle("/userschedule/{uid:[a-zA-Z0-9]+}/",
			M(initializeAddUserScheduleHandler(logConf, usConf, tConf), auth)).
			Methods(POST)

		// Party server
		s.Handle("/party/review/member",
			M(initializePartyReviewMemberHandler(logConf, pConf, uConf, tConf), auth)).
			Methods(POST)
		s.Handle("/party/review/done/{reviewer:[a-zA-Z0-9]+}",
			M(initializePartyReviewMemberDoneHandler(logConf, pConf, uConf, tConf), auth)).
			Methods(GET)
		s.Handle("/party/{uid:[a-zA-Z0-9]+}/{beginDateTime}/{endDateTime}",
			M(initializePartyHandler(logConf, pConf, uConf, tConf), auth)).
			Methods(GET)

		// Tag server
		// [Note] Longer path should be upper side
		s.Handle("/tags/{ttid:[0-9]+}",
			M(initializeTagsHandler(logConf, tConf), auth)).
			Methods(GET)
		s.Handle("/tags",
			M(initializeTagsHandler(logConf, tConf), auth)).
			Methods(GET)

		// User server
		s.Handle("/user/public/{uid:[a-zA-Z0-9]+}",
			M(initializeUserPublicHandler(logConf, uConf, tConf), auth)).
			Methods(GET)
		s.Handle("/user/register",
			M(initializeUserRegisterHandler(logConf, uConf, tConf), auth)).
			Methods(POST)
		s.Handle("/user/{uid:[a-zA-Z0-9]+}",
			M(initializeUserHandler(logConf, uConf, tConf), auth)).
			Methods(GET)
		// Block list
		s.Handle("/user/block",
			M(initializeUserBlockRegisterHandler(logConf, uConf, tConf), auth)).
			Methods(POST)

		// Health check
		s.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			db, err := sql.Open("mysql", tConf.DSN)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{
					"ok":      "false",
					"version": fmt.Sprintf("higher v%s", Version),
				})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			db.Close()

			json.NewEncoder(w).Encode(map[string]string{
				"ok":      "true",
				"version": fmt.Sprintf("higher v%s", Version),
			})
			w.WriteHeader(http.StatusOK)
		}).Methods(GET)

		// Launch web server
		log.Println("http:", *httpAddr)
		errChan <- http.ListenAndServe(*httpAddr, r)
	}()

	log.Fatal(<-errChan)
}
