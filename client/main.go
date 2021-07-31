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
	cli := client.New(setup.New())
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		timer := time.NewTimer(time.Second * 60)
		outbox := cli.Outbox()
		for {
			select {
			case bytes := <-outbox:
				log.Printf("Sever message: %s", bytes)
			case <-ticker.C:
				cli.Send(append([]byte{0x00, 0x08}, []byte("#foo#bar")...))
			case <-timer.C:
				cli.Stop()
			}
		}
	}()

	if err := cli.Run(); err != nil {
		log.Print("main: ", err)
	}
}
