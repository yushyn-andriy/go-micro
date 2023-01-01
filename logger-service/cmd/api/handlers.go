package main

import (
	"errors"
	"net/http"

	"github.com/yushyn-andriy/authentication/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.readJson(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.writeError(w, errors.New("Cannot insert log entry."), http.StatusInternalServerError)
		return
	}
	app.writeJson(w, http.StatusAccepted, jsonResponse{Error: false, Message: "logged"})
}
