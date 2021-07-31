package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	if err := godotenv.Load(); err != nil {
		log.Print("WARNING ", errors.Wrap(err, "init"))
		log.Print("Loading default setup")
	}
}

func main() {
	log.Println("Client run")
}
