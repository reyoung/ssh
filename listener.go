package ssh

import (
	"fmt"
	"io"
	"net"
)

type AcceptResult struct {
	RWC           io.ReadWriteCloser
	RemoteAddress string
	RemotePort    int
}

type RemoteListener interface {
	Accept() (*AcceptResult, error)
	Close() error
	Port() int
}

type stdListener struct {
	listener net.Listener
}

func (s *stdListener) Accept() (*AcceptResult, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("accept failed: %w", err)
	}
	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
	return &AcceptResult{
		RWC:           conn,
		RemoteAddress: tcpAddr.IP.String(),
		RemotePort:    tcpAddr.Port,
	}, nil
}

func (s *stdListener) Close() error {
	return s.listener.Close()
}

func (s *stdListener) Port() int {
	addr := s.listener.Addr()
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		return tcpAddr.Port
	}
	return -1
}

func defaultListen(address string) (RemoteListener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen failed: %w", err)
	}

	return &stdListener{listener: listener}, nil
}
