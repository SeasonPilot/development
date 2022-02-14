package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func main() {
	cfg := api.DefaultConfig()
	cfg.Address = "172.19.30.30:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      "redis",
		Name:    "redis",
		Tags:    []string{"primary", "v1"},
		Port:    8000,
		Address: "127.0.0.1",
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://127.0.0.1:8021/health",
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "redis"`)
	if err != nil {
		panic(err)
	}
	for key := range data {
		fmt.Println(key)
	}
}
