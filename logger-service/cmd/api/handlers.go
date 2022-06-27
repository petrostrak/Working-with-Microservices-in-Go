package main

import (
	"logger/data"
	"net/http"

	"github.com/petrostrak/toolbox"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var (
	tool toolbox.Tools
)

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into a var
	var requestPayload JSONPayload
	_ = tool.ReadJSON(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		tool.ErrorJSON(w, err)
		return
	}

	resp := toolbox.JSONResponse{
		Error:   false,
		Message: "logged",
	}

	tool.WriteJSON(w, http.StatusAccepted, resp)
}
