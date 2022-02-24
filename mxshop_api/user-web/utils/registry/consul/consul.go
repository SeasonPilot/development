package consul

import (
	"fmt"

	"mxshop-api/user-web/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type RegisterClient interface {
	Register(srvID, name string, tags []string, port int, addr string) error
	Deregister(srvID string) error
}

type Client struct {
	Client *api.Client
}

// NewConsulClient 函数签名返回接口类型
func NewConsulClient(host string, port int) RegisterClient {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &Client{
		Client: client,
	}
}

// Register 注册服务到 consul, 即服务注册
func (c Client) Register(srvID, name string, tags []string, port int, addr string) error {
	err := c.Client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      srvID,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: addr,
		Check: &api.AgentServiceCheck{
			//   注意端口是 consul 的端口！！！！ 和路径，协议是 http
			HTTP:                           fmt.Sprintf("http://%s:%d/v1/health/checks/%s", addr, global.SrvConfig.ConsulInfo.Port, global.SrvConfig.Name),
			Interval:                       "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func (c Client) Deregister(srvID string) error {
	err := c.Client.Agent().ServiceDeregister(srvID)
	if err != nil {
		zap.S().Errorf("注销服务失败:%s %s", global.SrvConfig.Name, srvID)
		return err
	}
	zap.S().Info("注销成功")
	return nil
}
