package hub

import (
	"fmt"
	"log"
	"net"
	"regexp"

	"tcpServer/server/client"
	"tcpServer/server/interfaces"

	"github.com/KnBrBz/message"

	"github.com/pkg/errors"
)

type H struct {
	host      string
	clients   map[string]*client.C
	reg       chan *client.C
	unreg     chan string
	broadcast chan *message.M
	tagNumber uint64
	done      chan struct{}
}

func New(stp interfaces.Setup) *H {
	return &H{
		host:      stp.Host(),
		clients:   make(map[string]*client.C),
		reg:       make(chan *client.C, capacity),
		unreg:     make(chan string, capacity),
		broadcast: make(chan *message.M, capacity),
		done:      make(chan struct{}),
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

	go h.eventsHandler()

	for {
		select {
		case <-h.done:
			return
		default:
			tcpConn, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Print(errors.Wrap(err, funcTitle))
				continue
			}
			h.reg <- client.New(tcpConn, h, h.nextTag(), tagReg)
		}
	}
}

func (h *H) Stop() {
	close(h.done)
}

func (h *H) Broadcast(msg *message.M) {
	h.broadcast <- msg
}

func (h *H) Unreg(uid string) {
	h.unreg <- uid
}

func (h *H) eventsHandler() {
	for {
		select {
		case c := <-h.reg:
			uid := c.UID()
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
		case <-h.done:
			return
		}
	}
}

func (h *H) nextTag() string {
	h.tagNumber++
	return fmt.Sprintf("#%d#", h.tagNumber)
}
