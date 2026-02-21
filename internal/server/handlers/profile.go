package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/glimesh/broadcast-box/internal/environment"
	"github.com/glimesh/broadcast-box/internal/server/authorization"
	"github.com/glimesh/broadcast-box/internal/server/helpers"
	"github.com/glimesh/broadcast-box/internal/webrtc"
	"github.com/glimesh/broadcast-box/internal/webrtc/sessions/manager"
)

func getProfileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		return
	}

	authUserHeader := os.Getenv(environment.AuthenticatedUserHeader)

	if authUserHeader == "" {
		helpers.LogHTTPError(
			responseWriter,
			"Authorization not enabled",
			http.StatusForbidden)
		return
	}

	streamKey := request.Header.Get(authUserHeader)

	if streamKey == "" {
		helpers.LogHTTPError(
			responseWriter,
			"No authorized user",
			http.StatusUnauthorized)
		return
	}

	getToken := func(streamKey string) (string, error) {
		token, _ := authorization.GetExistingProfileToken(streamKey)
		if token != "" {
			return token, nil
		}
		token, err := authorization.CreateProfile(streamKey)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	token, err := getToken(streamKey)

	if err != nil {
		return
	}

	responseWriter.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(responseWriter).Encode(map[string]any{"token": token}); err != nil {
		helpers.LogHTTPError(
			responseWriter,
			"Internal Server Error",
			http.StatusInternalServerError)
	}
}

func resetProfileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	authUserHeader := os.Getenv(environment.AuthenticatedUserHeader)

	if authUserHeader == "" {
		helpers.LogHTTPError(
			responseWriter,
			"Authorization not enabled",
			http.StatusForbidden)
		return
	}

	streamKey := request.Header.Get(authUserHeader)

	if streamKey == "" {
		return
	}

	err := authorization.ResetProfileToken(streamKey)

	if err != nil {
		log.Println("API.Log: Error resetting profile", err)
		helpers.LogHTTPError(responseWriter, "Could not reset token", http.StatusInternalServerError)
	}

	session, found := manager.SessionsManager.GetSessionByID(streamKey)
	if found == true {
		webrtc.HandleWHIPDelete(session.Host.Load().ID)
	}
	responseWriter.Header().Add("Content-Type", "text/plain")
	responseWriter.Write([]byte("OK\n"))
}
