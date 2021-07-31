package client

import (
	"bufio"
	"log"
	"net"
)

type C struct {
	uid string
	// tag  string
	conn *net.TCPConn
}

func New(conn *net.TCPConn) *C {
	return &C{
		uid:  conn.RemoteAddr().String(),
		conn: conn,
	}
}

func (c *C) Run() {
	go c.run()
}

func (c *C) run() {
	ipStr := c.conn.RemoteAddr().String()
	defer func() {
		log.Println("disconnected :" + ipStr)
		c.conn.Close()
	}()
	var bytes []byte
	reader := bufio.NewReader(c.conn)

	for {
		_, err := reader.Read(bytes)
		if err != nil {
			return
		}
		log.Println(c.uid + ":" + string(bytes))
		// Here the message is changed to broadcast
		//boradcastMessage(conn.RemoteAddr().String() + ":" + string(message))
	}
}

func (c *C) Uid() string {
	return c.uid
}
