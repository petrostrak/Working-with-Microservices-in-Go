package main

import (
	"net/http"

	"github.com/petrostrak/toolbox"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var tool toolbox.Tools

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tool.WriteJSON(w, http.StatusOK, payload)
}
