package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gitee.com/nelsonjs/go-mini.git/app"
	"gitee.com/nelsonjs/go-mini.git/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
)

type Client struct {
	cli      *clientv3.Client
	config   *config.Config
	Services []config.ServiceEnv
}

func GetClient(etcd EtcdManager, config *config.Config) *Client {
	c := &Client{
		cli:    etcd.client(),
		config: config,
	}
	fixClient(c, config)
	return c
}

func fixClient(c *Client, config *config.Config) {
	for _, v := range config.RSC.Services {
		srv, ok := app.Application().Services()[v.ServiceName]
		if !ok {
			continue
		}
		c.register(v.ServiceName, srv.Ctor, srv.IfacePtr)
	}
}

func (c *Client) register(serviceName string, ctor, ifacePtr interface{}) error {
	etcdResolver, err := resolver.NewBuilder(c.cli)
	if err != nil {
		panic(err)
	}
	confKey, err := filterDial(c.config.EC.ServiceNamePrefix, serviceName, c.config.EC.Version, c.cli)
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(fmt.Sprintf("etcd:///%s", confKey), grpc.WithResolvers(etcdResolver), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("grpc.Dial error: %v\n", err)
		return err
	}
	f := func(cc grpc.ClientConnInterface) interface{} {
		ret := reflect.ValueOf(ctor).Call([]reflect.Value{reflect.ValueOf(cc)})
		return ret[0].Interface()
	}
	app.GetContainer().Register(f(conn), ifacePtr)
	return nil
}

func filterDial(prefix, key, version string, cli *clientv3.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, fmt.Sprintf("%s/%s", prefix, key), clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("cli.Get fail,the error: %v", err)

	}
	for _, v := range resp.Kvs {
		conf := &Conf{}
		json.Unmarshal(v.Value, conf)
		if conf.Metadata.Version == version {
			return string(v.Key), nil
		}
	}
	return "", fmt.Errorf("service: %s/%s is not exists with version: %s", prefix, key, version)
}
