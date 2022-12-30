package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		false, "Hit the broker", nil,
	}
	_ = app.writeJson(w, http.StatusOK, payload)
}
