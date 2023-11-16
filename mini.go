package mini

import (
	"net"

	"gitee.com/nelsonjs/go-mini.git/config"
	"google.golang.org/grpc"
)

type Service interface {
	// Server() *grpc.Server
	WithGrpc(fn func(srv *grpc.Server)) Service
	WithHttp(fn func(l net.Listener) error) Service
	Run()
}

func NewService(name string, config *config.Config) Service {
	return newService(name, config)
}
