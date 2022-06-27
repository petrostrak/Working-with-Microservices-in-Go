package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/petrostrak/toolbox"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var tool toolbox.Tools

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := tool.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tool.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the DB
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		tool.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		tool.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s\n", user.Email),
		Data:    user,
	}

	tool.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
