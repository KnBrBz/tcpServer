package client

import (
	"bufio"
	"log"
	"net"
	"tcpServer/client/interfaces"
	"tcpServer/message"

	"github.com/pkg/errors"
)

type C struct {
	inbox   chan []byte
	outbox  chan []byte
	srvHost string
	host    string
	done    chan struct{}
}

func New(stp interfaces.Setup) *C {
	return &C{
		inbox:   make(chan []byte),
		srvHost: stp.ServerHost(),
		host:    stp.Host(),
		done:    make(chan struct{}),
	}
}

func (c *C) Run() (err error) {
	const funcTitle = packageTitle + "*C.Run"
	var tcpAddr, laddr *net.TCPAddr

	if tcpAddr, err = net.ResolveTCPAddr("tcp", c.srvHost); err != nil {
		return errors.Wrapf(err, "%s: server host", funcTitle)
	}
	if laddr, err = net.ResolveTCPAddr("tcp", c.host); err != nil {
		return errors.Wrapf(err, "%s: host", funcTitle)
	}
	conn, err := net.DialTCP("tcp", laddr, tcpAddr)
	if err != nil {
		return errors.Wrap(err, funcTitle)
	}
	defer conn.Close()

	go c.read(conn)
	// console chat feature joins
	for {
		select {
		case bytes := <-c.inbox:
			msg := message.New(message.HeadLength, bytes)
			if err := msg.Validate(nil); err != nil {
				log.Printf("Inbox message `%s` not valid: %v", msg.Body(), err)
			}
			if _, err := conn.Write(msg.Body()); err != nil {
				return errors.Wrap(err, funcTitle)
			}
		case <-c.done:
			return
		}
	}
}

func (c *C) Outbox() <-chan []byte {
	return c.outbox
}

func (c *C) Send(bytes []byte) {
	c.inbox <- bytes
}

func (c *C) Stop() {
	close(c.done)
}

func (c *C) read(conn *net.TCPConn) {
	const funcTitle = packageTitle + "*C.Read"
	var bytes []byte
	reader := bufio.NewReader(conn)
	for {
		_, err := reader.Read(bytes)
		if err != nil {
			log.Print(errors.Wrap(err, funcTitle))
			break
		}
		msg := message.New(message.HeadLength, bytes)
		if err := msg.Validate(nil); err != nil {
			log.Printf("Outbox message `%s` not valid: %v", msg.Body(), err)
			continue
		}
		c.outbox <- bytes
	}
}