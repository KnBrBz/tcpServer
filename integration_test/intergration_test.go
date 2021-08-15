package intagrationtest

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"tcpServer/client/client"
	"tcpServer/integration_test/fakesetup"
	"tcpServer/server/hub"

	"github.com/pkg/errors"
)

func TestClientsServer(t *testing.T) {
	serverPort := 9999
	mainHost := "127.0.0.1:"
	serverHost := mainHost + strconv.Itoa(serverPort)
	maxClients := 10
	errChan := make(chan error, maxClients+1)
	// init server
	srv := hub.New(fakesetup.NewServers(serverHost))
	defer srv.Stop()

	go func() {
		if err := srv.Run(); err != nil {
			errChan <- err
		}
	}()
	time.Sleep(time.Second) // TODO make some check that server is running
	// init 10 clients
	clis := make([]*client.C, maxClients)
	defer func() {
		for _, cli := range clis {
			if cli != nil {
				cli.Stop()
			}
		}
	}()

	for i := 0; i < maxClients; i++ {
		clientHost := mainHost + strconv.Itoa(serverPort-i-1)
		cli := client.New(fakesetup.NewClients(clientHost, serverHost))
		clis[i] = cli

		go func() {
			if err := cli.Run(); err != nil {
				errChan <- err
			}
		}()
	}
	time.Sleep(time.Second) // TODO make some check that server clients are running
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
	// send message
	msg := append([]byte{0x00, 0x08}, []byte("bar#foo#")...)

	var wg sync.WaitGroup

	wg.Add(len(clis))

	for i, cli := range clis {
		go func(i int, cli *client.C) {
			defer wg.Done()

			outbox := cli.Outbox()
			timer := time.NewTimer(time.Second)

			defer timer.Stop()
			select {
			case cliMsg := <-outbox:
				if !assertBytesEqual(cliMsg, msg) {
					errChan <- errors.Errorf("client `%d`: expected message `%s`, got `%s`", i, msg, cliMsg)
				}
			case <-timer.C:
				errChan <- errors.Errorf("client `%d`: timeout", i)
			}
		}(i, cli)
	}

	i := rand.Intn(10)
	clis[i].Send(msg)
	wg.Wait()
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
	// send tagged message
	srcI := rand.Intn(10)
	dstTag := intToTag(rand.Intn(10) + 1)
	content := []byte(dstTag + "foobar")
	msg = append(int16ToBytes(int16(len(content))), content...)
	flag := &boolFlag{}

	wg.Add(len(clis))

	for i, cli := range clis {
		go func(i int, cli *client.C) {
			defer wg.Done()

			outbox := cli.Outbox()
			timer := time.NewTimer(time.Second)

			defer timer.Stop()
			select {
			case cliMsg := <-outbox:
				switch {
				case flag.Read():
					errChan <- errors.Errorf("client `%d`: message `%s` was already received", i, cliMsg)
				default:
					if !assertBytesEqual(cliMsg, msg) {
						errChan <- errors.Errorf("client `%d`: expected message `%s`, got `%s`", i, msg, cliMsg)
					}

					flag.Set(true)
				}
			case <-timer.C:
				if !flag.Read() {
					errChan <- errors.Errorf("client `%d`: timeout", i)
				}
			}
		}(i, cli)
	}

	clis[srcI].Send(msg) // TODO it is possible to send message to myself, because we don't know clients tag
	wg.Wait()
	select {
	case err := <-errChan:
		t.Fatal(err)
	default:
	}
}

type boolFlag struct {
	value bool
	mux   sync.RWMutex
}

func (bf *boolFlag) Set(value bool) {
	bf.mux.Lock()
	defer bf.mux.Unlock()
	bf.value = value
}

func (bf *boolFlag) Read() bool {
	bf.mux.RLock()
	defer bf.mux.RUnlock()

	return bf.value
}

func int16ToBytes(i int16) []byte {
	var h, l uint8 = uint8(i >> 8), uint8(i & 0xff)
	return []byte{h, l}
}

func intToTag(i int) string {
	return "#" + strconv.Itoa(i) + "#"
}

func assertBytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
