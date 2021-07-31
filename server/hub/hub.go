package hub

import (
	"log"
	"net"
	"tcpServer/server/client"
	"tcpServer/setup"

	"github.com/pkg/errors"
)

type H struct {
	host    string
	clients map[string]*client.C
}

func New(stp *setup.S) *H {
	return &H{
		host:    stp.Host(),
		clients: make(map[string]*client.C),
	}
}

func (h *H) Run() (err error) {
	const funcTitle = packageTitle + "*H.Run"
	tcpAddr, err := net.ResolveTCPAddr("tcp", h.host)
	if err != nil {
		return errors.Wrap(err, funcTitle)
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errors.Wrap(err, funcTitle)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Print(errors.Wrap(err, funcTitle))
			continue
		}
		c := client.New(tcpConn)
		uid := c.Uid()
		log.Println("A client connected : " + uid)
		h.clients[uid] = c
		c.Run()
	}
}
