package service

import "testing"
import nn "github.com/op/go-nanomsg"
import "github.com/stretchr/testify/assert"

func TestEcho(t *testing.T) {
	serv, err := New("tcp://127.0.0.1:9083")
	assert.NoError(t, err)
	go serv.echoServe()
	go serv.echoServe()

	testEcho := func(msg []byte, iteration int) {
		sock, err := nn.NewReqSocket()
		defer sock.Close()
		for i := 0; i < iteration; i++ {
			assert.NoError(t, err)
			_, err = sock.Connect("tcp://127.0.0.1:9083")
			assert.NoError(t, err)
			_, err = sock.Send(msg, len(msg))
			assert.NoError(t, err)
			buf, err := sock.Recv(0)
			assert.NoError(t, err)
			assert.Equal(t, msg, buf)
		}
	}

	go testEcho([]byte("abc"), 10000)
	go testEcho([]byte("bcd"), 10000)
}
