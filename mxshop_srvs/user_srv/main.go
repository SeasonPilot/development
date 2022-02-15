package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"net"

	"mxshop-srvs/user_srv/global"
	"mxshop-srvs/user_srv/handler"
	"mxshop-srvs/user_srv/initialization"
	"mxshop-srvs/user_srv/model"
	"mxshop-srvs/user_srv/proto"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitDB()

	//密码加密
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	pw := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(pw)

	for i := 0; i < 10; i++ {
		user := model.User{
			Mobile:   fmt.Sprintf("1879876789%d", i),
			Password: pw,
			NickName: fmt.Sprintf("bobby%d", i),
		}
		global.DB.Save(&user)
	}

	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	zap.S().Info(*ip, *port)

	// 注册服务
	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handler.UserServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}

	// 注册 grpc 服务健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	// 注册服务到 consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      global.ServiceConfig.Name,
		Name:    global.ServiceConfig.Name,
		Tags:    []string{"primary", "v1"},
		Port:    *port,
		Address: "10.94.62.100",
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d", "10.94.62.100", *port),
			DeregisterCriticalServiceAfter: "15s",
		},
	})
	if err != nil {
		panic(err)
	}

	// 启动服务
	err = g.Serve(listen)
	if err != nil {
		panic(err)
	}
}
