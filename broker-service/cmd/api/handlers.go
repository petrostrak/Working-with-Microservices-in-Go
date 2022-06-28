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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tool.WriteJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
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
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		tool.ErrorJSON(w, errors.New("unknown action"))
	}
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
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
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		tool.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		tool.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService toolbox.JSONResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		tool.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	tool.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	// call the service
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		tool.ErrorJSON(w, errors.New("error calling logger service"))
		return
	}

	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "logged!"

	tool.WriteJSON(w, http.StatusAccepted, payload)

}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	// call the service - post to mail service
	request, err := http.NewRequest("POST", "http://mailer-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		tool.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "mail sent to " + msg.To

	tool.WriteJSON(w, http.StatusAccepted, payload)
}
