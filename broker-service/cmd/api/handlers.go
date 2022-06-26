package main

import (
	"net/http"

	"github.com/petrostrak/toolbox"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var tool toolbox.Tools

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tool.WriteJSON(w, http.StatusOK, payload)
}
