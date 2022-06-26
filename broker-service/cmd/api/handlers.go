package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/petrostrak/toolbox"
)

var (
	tool toolbox.Tools
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tool.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tool.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		tool.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json to send to auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	// call the service
	request, err := http.NewRequest("POST", "http://auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}
	defer request.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		tool.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		tool.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable to read response.Body into
	var jsonFromService toolbox.JSONResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		tool.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// here we have a valid login
	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	tool.WriteJSON(w, http.StatusAccepted, payload)
}
