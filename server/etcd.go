package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitee.com/nelsonjs/go-mini.git/utils.go"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type EtcdManager interface {
	Register(key string, addrs ...string)
	UnRegister(key string)
	client() *clientv3.Client
}

type etcdService struct {
	Cli    *clientv3.Client
	em     endpoints.Manager
	close  chan bool
	config *EtcdConfig
}

type EtcdConfig struct {
	EndPoints         []string
	Timeout           int
	Port              string
	ServiceNamePrefix string
	Version           string
}

type Conf struct {
	Addr     string
	Metadata Metadata
}

type Metadata struct {
	Version string
}

func NewEtcd(config *EtcdConfig) EtcdManager {
	if len(config.EndPoints) == 0 {
		panic(errors.New("ETCD endpoint is empty"))
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 5
	}
	cli := newEtcdClient(config.EndPoints, timeout)
	return &etcdService{
		Cli:    cli,
		em:     newEtcdEndpointsManager(cli),
		close:  make(chan bool),
		config: config,
	}
}

func newEtcdEndpointsManager(cli *clientv3.Client) endpoints.Manager {
	em, err := endpoints.NewManager(cli, viper.GetString("etcd.key_prefix"))
	if err != nil {
		panic(err)
	}
	return em
}

func newEtcdClient(endPoints []string, timeout int) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return cli
}

func (s *etcdService) Register(key string, addrs ...string) {
	lease, err := s.Cli.Grant(context.TODO(), 7)
	if err != nil {
		panic(err)
	}
	addr := fmt.Sprintf("%s:%s", getIp(), s.config.Port)
	target := fmt.Sprintf("%s/%s/%s", s.config.ServiceNamePrefix, key, addr)
	s.em.AddEndpoint(context.TODO(), target, endpoints.Endpoint{
		Addr: addr,
		Metadata: Metadata{
			Version: s.config.Version,
		},
	}, clientv3.WithLease(lease.ID))
	s.keepAlive(lease.ID)
}

func (s *etcdService) UnRegister(key string) {
	s.close <- true
	s.em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s%s", s.config.ServiceNamePrefix, key))
	s.Cli.Close()
}

func (s *etcdService) client() *clientv3.Client {
	return s.Cli
}

func getIp() string {
	addr, err := utils.GetOutboundIP()
	if err != nil {
		addr, err = utils.GuessExternalIP()
	}
	if err != nil {
		panic(err)
	}
	return addr.String()
}

func (s *etcdService) keepAlive(id clientv3.LeaseID) {
	dur := time.Duration(time.Second * time.Duration(5))
	timer := time.NewTicker(dur)

	go func(leaseId clientv3.LeaseID) {
		for {
			select {
			case <-timer.C:
				_, err := s.Cli.KeepAliveOnce(context.TODO(), leaseId)
				if err != nil {
					fmt.Printf("ETCD KeepAlive error: %v", err)
				}
			case <-s.close:
				goto EXIT
			}
		}
	EXIT:
		fmt.Println("KeepAlive Exit")
	}(id)
}
