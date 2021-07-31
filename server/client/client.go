package client

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"tcpServer/message"
	"tcpServer/server/client/requisites"
	"tcpServer/server/interfaces"

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
	go c.read()
	go c.write()
}

func (c *C) Send(msg *message.M) {
	c.inbox <- msg
}

func (c *C) read() {
	const funcTitle = packageTitle + "*C.read"
	ipStr := c.conn.RemoteAddr().String()
	defer func() {
		log.Println("disconnected :" + ipStr)
		c.conn.Close()
		c.hub.Unreg(c.req.UID)
	}()
	var bytes []byte
	reader := bufio.NewReader(c.conn)
	for {
		_, err := reader.Read(bytes)
		if err != nil {
			return
		}
		msg := message.New(msgHeadLength, bytes)
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

func (c *C) write() {
	const funcTitle = packageTitle + "write"
	for msg := range c.inbox {
		if tag := msg.Tag(); len(tag) > 0 && tag != c.req.Tag {
			return
		}
		if _, err := c.conn.Write(msg.Body()); err != nil {
			log.Println(errors.Wrap(err, funcTitle))
		}
	}
}
