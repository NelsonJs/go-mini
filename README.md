# go-mini

#### 介绍
一款简单的的微服务框架

#### 软件架构
采用Etcd做服务发现功能，集成grpc功能。


#### 安装教程

```
go get gitee.com/nelsonjs/go-mini.git
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


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
