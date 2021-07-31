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
	reg     chan *client.C
	unreg   chan string
}

func New(stp *setup.S) *H {
	return &H{
		host:    stp.Host(),
		clients: make(map[string]*client.C),
		reg:     make(chan *client.C, 10),
		unreg:   make(chan string, 10),
	}
}

func (h *H) eventsHandler(done <-chan struct{}) {
	for {
		select {
		case c := <-h.reg:
			uid := c.Uid()
			log.Println("A client connected : " + uid)
			h.clients[uid] = c
			c.Run()
		case uid := <-h.unreg:
			log.Println("A client disconnected : " + uid)
			delete(h.clients, uid)
		case <-done:
			return
		}
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
	done := make(chan struct{})
	go h.eventsHandler(done)
	defer close(done)
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Print(errors.Wrap(err, funcTitle))
			continue
		}
		h.reg <- client.New(tcpConn, h.unreg)
	}
}
