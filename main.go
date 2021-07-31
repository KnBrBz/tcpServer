package main

import (
	"bufio"
	"log"
	"net"
	"tcpServer/setup"

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

// used to record all client connections
var ConnMap map[string]*net.TCPConn

func main() {
	var tcpAddr *net.TCPAddr
	ConnMap = make(map[string]*net.TCPConn)
	stp := setup.New()
	tcpAddr, _ = net.ResolveTCPAddr("tcp", stp.Host())

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		log.Println("A client connected : " + tcpConn.RemoteAddr().String())
		// new connection to join map
		ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
		go tcpPipe(tcpConn)
	}
}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		log.Println("disconnected :" + ipStr)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		log.Println(conn.RemoteAddr().String() + ":" + string(message))
		// Here the message is changed to broadcast
		boradcastMessage(conn.RemoteAddr().String() + ":" + string(message))
	}
}

func boradcastMessage(message string) {
	b := []byte(message)
	// traverse all clients and send messages
	for _, conn := range ConnMap {
		conn.Write(b)
	}
}
