package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PORT = "80"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", PORT)

	// define http server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	// start the server
	if err = srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var (
		count      int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			count++
		} else {
			log.Println("Connected to RabbitMQ..")
			connection = c
			break
		}

		if count > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Backing of...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
