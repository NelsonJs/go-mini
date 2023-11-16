package server

import (
	"fmt"
	"net"

	"github.com/soheilhy/cmux"
)

type Mux interface {
	GrpcNetListener() net.Listener
	HttpNetListener() net.Listener
	Serve() error
}

type mux struct {
	m cmux.CMux
}

func GetMux(port string) Mux {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	return &mux{
		m: cmux.New(l),
	}
}

func (m *mux) GrpcNetListener() net.Listener {
	return m.m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
}

func (m *mux) HttpNetListener() net.Listener {
	return m.m.Match(cmux.HTTP1Fast())
}

func (m *mux) Serve() error {
	return m.m.Serve()
}
