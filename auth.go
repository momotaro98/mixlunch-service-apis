package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/momotaro98/mixlunch-service-api/logger"
)

func AuthMiddle(activate bool, loggerConfig *logger.Config, appCredentialFilePath string) MFunc {
	if !activate {
		// Do nothing and just pass to next http handler
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	opt := option.WithCredentialsFile(appCredentialFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return &authHandler{
			next:   next,
			logger: logger.ProvideLogger(loggerConfig),
			app:    app,
		}
	}
}

type authHandler struct {
	next   http.Handler
	logger logger.Logger
	app    *firebase.App
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqId := r.Header.Get(XRequestId)

	client, err := h.app.Auth(ctx)
	if err != nil {
		h.logger.Log(logger.Error, reqId, fmt.Sprintf("Firebase app.Auth error. Error: %+v\n", err.Error()))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Bearer" {
		h.logger.Log(logger.Warn, reqId, fmt.Sprintf("Request header Authorization is nothing or invalid.\n"))
		http.Error(w, "Authorization header is not valid.", http.StatusBadRequest)
		return
	}
	token := auth[1]

	if _, err := client.VerifySessionCookie(ctx, token); err != nil { // VerifySessionCookie is too slow because it connects to Firebase online
		h.logger.Log(logger.Warn, reqId, fmt.Sprintf("error verifying token: %+v\n", err.Error()))
		http.Error(w, "Token is not valid.", http.StatusUnauthorized)
		return
	}

	h.next.ServeHTTP(w, r)
}
