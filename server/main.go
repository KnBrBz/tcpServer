package main

import (
	"log"
	"tcpServer/server/hub"
	"tcpServer/setup"

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
	if err := hub.New(
		setup.New(),
	).Run(); err != nil {
		log.Print(errors.Wrap(err, "main"))
	}
}