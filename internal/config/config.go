package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     int
	DBFile   string
	Password string
}

func Load() (Config, error) {
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		port = "7540"
	}
	addr, err := strconv.Atoi(port)
	if err != nil {
		return Config{}, fmt.Errorf("port must be a number")
	}
	dbFile, ok := os.LookupEnv("TODO_DBFILE")
	if !ok {
		dbFile = "./pkg/db/scheduler.db"
	}
	password := os.Getenv("TODO_PASSWORD")

	return Config{Port: addr, DBFile: dbFile, Password: password}, nil
}
