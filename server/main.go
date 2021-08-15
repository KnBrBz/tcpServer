package main

import (
	"log"
	"tcpServer/server/hub"
	"tcpServer/server/setup"

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
	defer func() {
		if err := recover(); err != nil {
			log.Print("main: ", err)
		}
	}()

	if err := hub.New(
		setup.New(),
	).Run(); err != nil {
		log.Print(errors.Wrap(err, "main"))
	}
}
