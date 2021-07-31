package hub

import (
	"fmt"
	"log"
	"net"
	"regexp"

	"tcpServer/message"
	"tcpServer/server/client"
	"tcpServer/server/setup"

	"github.com/pkg/errors"
)

type H struct {
	host      string
	clients   map[string]*client.C
	reg       chan *client.C
	unreg     chan string
	broadcast chan *message.M
	tagNumber uint64
}

func New(stp *setup.S) *H {
	return &H{
		host:      stp.Host(),
		clients:   make(map[string]*client.C),
		reg:       make(chan *client.C, 10),
		unreg:     make(chan string, 10),
		broadcast: make(chan *message.M, 10),
	}
}

func (h *H) Run() (err error) {
	const funcTitle = packageTitle + "*H.Run"
	tagReg := regexp.MustCompile("#.+#")
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
	go h.eventsHandler(done, tagReg)
	defer close(done)
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Print(errors.Wrap(err, funcTitle))
			continue
		}
		h.reg <- client.New(tcpConn, h, h.nextTag(), tagReg)
	}
}

func (h *H) Broadcast(msg *message.M) {
	h.broadcast <- msg
}

func (h *H) Unreg(uid string) {
	h.unreg <- uid
}

func (h *H) eventsHandler(done <-chan struct{}, tagReg *regexp.Regexp) {
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
		case msg := <-h.broadcast:
			for _, c := range h.clients {
				c.Write(msg)
			}
		case <-done:
			return
		}
	}
}

func (h *H) nextTag() string {
	h.tagNumber++
	return fmt.Sprintf("#%d#", h.tagNumber)
}
