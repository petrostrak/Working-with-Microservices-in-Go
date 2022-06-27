package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	PORT = "80"
)

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", PORT)

	// define http server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	// start the server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
