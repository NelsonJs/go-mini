# go-mini

#### 介绍
一款简单的的微服务框架

#### 软件架构
采用Etcd做服务发现功能，集成grpc功能。


#### 安装教程

```
go get github.com/nelsonjs/go-mini.git
```
### 使用

#### 注册服务
与proto编译后的文件同一个目录，创建一个已.go结尾的文件，在init函数中调用注册函数
```
app.Application().Register("api", NewGreeterClient, (*GreeterClient)(nil))
```

#### 服务端初始化
```
service := mini.NewService("api", &config.Config{
		EC: config.EtcdConfig{
			Endpoints:         []string{"http://172.31.190.27:2379"},
			Timeout:           5,
			ServiceNamePrefix: "/test",
			Version:           "dev",
		},
		HC: config.HttpConfig{
			Port: 50051,
		},
	})
	service.WithGrpc(func(srv *grpc.Server) {
		helloworld.RegisterGreeterServer(srv, &server{})
	}).WithHttp(func(l net.Listener) error {
		i := iris.New()
		i.Get("/", func(ctx iris.Context) {
			ctx.JSON(map[string]string{"name": "p"})
		})
		return i.Run(iris.Listener(l))
	})

	service.Run()
```
#### 客户端调用grpc服务
```
mini.NewService("", &config.Config{
		EC: config.EtcdConfig{
			Endpoints:         []string{"http://172.31.190.27:2379"},
			ServiceNamePrefix: "/test",
			Version:           "dev",
		},
		RSC: config.ReferenceServiceConfig{
			Services: []config.ServiceEnv{
				{ServiceName: "api"},
			},
		},
	})
	app.GetContainer().Invoke(func(greeter helloworld.GreeterClient) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		r, err := greeter.SayHello(ctx, &helloworld.HelloRequest{
			Name: "111",
		})
		cancel()
		if err != nil {
			fmt.Printf("could not greet: %v", err)
		}
		fmt.Printf("Greeting: %s", r.GetMessage())
	})
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
