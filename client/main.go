package main

import (
	"log"
	"tcpServer/client/client"
	"tcpServer/client/setup"
	"time"

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
	const (
		sendInterval = time.Second * 5
		timeout      = time.Second * 60
	)

	cli := client.New(setup.New())

	go func() {
		ticker := time.NewTicker(sendInterval)
		timer := time.NewTimer(timeout)

		defer ticker.Stop()
		defer timer.Stop()

		outbox := cli.Outbox()

		for {
			select {
			case bytes := <-outbox:
				log.Printf("Server message: %s", bytes)
			case <-ticker.C:
				cli.Send(append([]byte{0x00, 0x08}, []byte("bar#foo#")...))
			case <-timer.C:
				cli.Stop()
			}
		}
	}()

	if err := cli.Run(); err != nil {
		log.Print("main: ", err)
	}
}
