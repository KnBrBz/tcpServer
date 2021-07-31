package client

import (
	"bufio"
	"log"
	"net"
	"tcpServer/server/interfaces"

	"github.com/pkg/errors"
)

type C struct {
	uid   string
	conn  *net.TCPConn
	inbox chan []byte
	hub   interfaces.Hub
}

func New(conn *net.TCPConn, hub interfaces.Hub) *C {
	return &C{
		uid:   conn.RemoteAddr().String(),
		conn:  conn,
		hub:   hub,
		inbox: make(chan []byte),
	}
}

func (c *C) Uid() string {
	return c.uid
}

func (c *C) Run() {
	go c.read()
	go c.write()
}

func (c *C) Send(bytes []byte) {
	c.inbox <- bytes
}

func (c *C) read() {
	ipStr := c.conn.RemoteAddr().String()
	defer func() {
		log.Println("disconnected :" + ipStr)
		c.conn.Close()
		c.hub.Unreg(c.uid)
	}()
	var bytes []byte
	reader := bufio.NewReader(c.conn)

	for {
		_, err := reader.Read(bytes)
		if err != nil {
			return
		}
		c.hub.Broadcast(bytes)
	}
}

func (c *C) Write(bytes []byte) {
	c.inbox <- bytes
}

func (c *C) write() {
	const funcTitle = packageTitle + "write"
	for bytes := range c.inbox {
		if _, err := c.conn.Write(bytes); err != nil {
			log.Println(errors.Wrap(err, funcTitle))
		}
	}
}
