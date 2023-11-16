package mini

import (
	"net"

	"gitee.com/nelsonjs/go-mini.git/config"
	"gitee.com/nelsonjs/go-mini.git/server"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
)

type service struct {
	name   string
	config *config.Config
	server *grpc.Server
	mux    server.Mux
	etcd   server.EtcdManager
}

func fixConfig(cnf *config.Config) {
	if cnf == nil {
		cnf = &config.Config{}
	}
	if cnf.EC.Timeout <= 0 {
		cnf.EC.Timeout = 5
	}
}

func newService(name string, config *config.Config) Service {
	srv := &service{
		name:   name,
		config: config,
		server: server.DefaultGrpcServer.GrpcServer(),
		mux:    server.GetMux(cast.ToString(config.HC.Port)),
		etcd: server.NewEtcd(&server.EtcdConfig{
			EndPoints:         config.EC.Endpoints,
			Timeout:           config.EC.Timeout,
			Port:              cast.ToString(config.HC.Port),
			ServiceNamePrefix: config.EC.ServiceNamePrefix,
			Version:           config.EC.Version,
		}),
	}
	fixConfig(config)
	server.GetClient(srv.etcd, config)
	return srv
}

func (s *service) WithGrpc(fn func(srv *grpc.Server)) Service {
	s.etcd.Register(s.name)
	fn(s.server)
	go func() {
		if err := s.server.Serve(s.mux.GrpcNetListener()); err != nil {
			panic(err)
		}
	}()
	return s
}

func (s *service) WithHttp(fn func(l net.Listener) error) Service {
	l := s.mux.HttpNetListener()
	go func(l net.Listener) {
		if err := fn(l); err != nil {
			panic(err)
		}
	}(l)
	return s
}

func (s *service) Client() *server.Client {
	return server.GetClient(s.etcd, s.config)
}

func (s *service) Run() {
	if err := s.mux.Serve(); err != nil {
		panic(err)
	}
}
