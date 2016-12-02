package service

import nn "github.com/op/go-nanomsg"

// Service for the parameter server.
type Service struct {
	rawSocket *nn.RepSocket
}

// New parameter server by addr
func New(addr string) (*Service, error) {
	sock, err := nn.NewRepSocket()
	if err != nil {
		return nil, err
	}
	_, err = sock.Bind(addr)
	if err != nil {
		return nil, err
	}
	return &Service{
		rawSocket: sock,
	}, nil
}

func (serv *Service) echoServe() error {
	for {
		buf, err := serv.rawSocket.Recv(0)
		if err != nil {
			return err
		}
		_, err = serv.rawSocket.Send(buf, len(buf))
		if err != nil {
			return err
		}
	}
}
