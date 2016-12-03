package service

import "github.com/go-mangos/mangos"

import "github.com/go-mangos/mangos/transport/tcp"

import "github.com/go-mangos/mangos/transport/inproc"

import "github.com/reyoung/GoPServer/param"

import "github.com/go-mangos/mangos/protocol/req"
import "github.com/go-mangos/mangos/protocol/rep"
import "fmt"

// Service for the parameter server.
type Service struct {
	params        *param.Parameters
	devSocketAddr string
	exposedSocket mangos.Socket
	devSocket     mangos.Socket
}

// New parameter server by addr
func New(addr string, sockname string) (*Service, error) {
	sock, err := rep.NewSocket()
	if err != nil {
		return nil, err
	}
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(addr); err != nil {
		return nil, err
	}
	devSocks, err := req.NewSocket()
	if err != nil {
		return nil, err
	}
	devSocks.AddTransport(inproc.NewTransport())
	devSocketAddr := fmt.Sprintf("inproc://%s", sockname)

	if err = devSocks.Listen(devSocketAddr); err != nil {
		return nil, err
	}
	go func() {
		mangos.Device(sock, devSocks)
	}()

	return &Service{devSocketAddr: devSocketAddr, exposedSocket: sock,
		devSocket: devSocks}, nil
}

func (serv *Service) echoServe() error {
	sock, err := rep.NewSocket()
	if err != nil {
		return err
	}
	sock.AddTransport(inproc.NewTransport())
	if err = sock.Dial(serv.devSocketAddr); err != nil {
		return err
	}
	for {
		msg, err := sock.Recv()
		if err != nil {
			return err
		}
		if err = sock.Send(msg); err != nil {
			return err
		}
	}
}

func (s *Service) Close() {
	s.devSocket.Close()
	s.exposedSocket.Close()
}
