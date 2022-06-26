package main

import (
	"authentication/data"
	"database/sql"
)

const (
	PORT = "8080"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}
