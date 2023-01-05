package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,ommitempty"`
	Log    LogPayload  `json:"log,ommitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		false, "Hit the broker", nil,
	}
	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.log(w, requestPayload.Log)
	default:
		app.writeError(w, errors.New("unknown action"))
	}

}

func (app *Config) log(w http.ResponseWriter, a LogPayload) {
	requestPayload, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Println(err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.writeError(w, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.writeError(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.writeError(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.writeError(w, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.writeError(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}
