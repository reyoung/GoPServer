package service

import "testing"

import "github.com/go-mangos/mangos/protocol/req"
import "github.com/stretchr/testify/assert"
import "sync"
import "github.com/go-mangos/mangos/transport/tcp"

func TestEcho(t *testing.T) {
	addr := "tcp://127.0.0.1:4321"
	serv, err := New(addr, "test_echo")
	assert.NoError(t, err)
	go serv.echoServe()
	// go serv.echoServe()
	var wg sync.WaitGroup
	wg.Add(2)
	testEcho := func(msg []byte, iteration int) {
		sock, err := req.NewSocket()
		if err != nil {
			panic(err)
		}
		sock.AddTransport(tcp.NewTransport())
		err = sock.Dial(addr)
		if err != nil {
			panic(err)
		}
		for i := 0; i < iteration; i++ {
			err = sock.Send(msg)
			if err != nil {
				panic(err)
			}
			buf, err := sock.Recv()
			if err != nil {
				panic(err)
			}
			assert.Equal(t, msg, buf)
		}
		wg.Done()
		assert.NoError(t, sock.Close())
	}

	go testEcho([]byte("abc"), 10000)
	go testEcho([]byte("bcd"), 10000)
	wg.Wait()
	serv.Close()
}
