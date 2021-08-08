package client

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"tcpServer/server/client/requisites"
	"tcpServer/server/interfaces"

	"github.com/KnBrBz/message"

	"github.com/pkg/errors"
)

type C struct {
	req    requisites.R
	conn   *net.TCPConn
	inbox  chan *message.M
	hub    interfaces.Hub
	tagReg *regexp.Regexp
}

func New(conn *net.TCPConn, hub interfaces.Hub, tag string, tagReg *regexp.Regexp) *C {
	return &C{
		conn:  conn,
		hub:   hub,
		inbox: make(chan *message.M),
		req: requisites.R{
			UID: conn.RemoteAddr().String(),
			Tag: tag,
		},
		tagReg: tagReg,
	}
}

func (c *C) Uid() string {
	return c.req.UID
}

func (c *C) Run() {
	done := make(chan struct{})
	go c.read(done)
	go c.write(done)
}

func (c *C) Send(msg *message.M) {
	c.inbox <- msg
}

func (c *C) read(done chan struct{}) {
	const funcTitle = packageTitle + "*C.read"
	ipStr := c.conn.RemoteAddr().String()
	defer func() {
		log.Println("disconnected: " + ipStr)
		c.hub.Unreg(c.req.UID)
		close(done)
		c.conn.Close()
	}()

	var bytes []byte = make([]byte, 0xffff+2)
	reader := bufio.NewReader(c.conn)
	for {
		n, err := reader.Read(bytes)
		if err != nil {
			log.Print(errors.Wrap(err, funcTitle))
			return
		}
		msg := message.New(msgHeadLength, bytes[:n])
		if err := msg.Validate(c.tagReg); err != nil {
			log.Print(errors.Wrapf(err, "%s, client %s", funcTitle, c.Uid()))
			continue
		}
		c.hub.Broadcast(msg)
	}
}

func (c *C) Write(msg *message.M) {
	c.inbox <- msg
}

func (c *C) write(done chan struct{}) {
	const funcTitle = packageTitle + "write"
	for {
		select {
		case msg := <-c.inbox:
			if tag := msg.Tag(); len(tag) > 0 && tag != c.req.Tag {
				log.Printf("skip tag `%s`, client tag `%s`", tag, c.req.Tag)
				continue
			}
			if _, err := c.conn.Write(msg.Body()); err != nil {
				log.Println(errors.Wrap(err, funcTitle))
			}
		case <-done:
			return
		}
	}
}
