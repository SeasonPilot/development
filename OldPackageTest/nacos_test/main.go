package main

import (
	"fmt"
	"strings"

	"development/OldPackageTest/nacos_test/config"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

func main() {
	// ServerConfig 是 slice
	sc := []constant.ServerConfig{
		{
			IpAddr: "172.19.30.30",
			Port:   8848,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         "731b0ec0-4df5-4b84-a2a4-667e41e933cc", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-srv.yaml",
		Group:  "dev",
	})
	if err != nil {
		panic(err)
	}

	// 从 io.Reader 读取 config
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	//err = viper.ReadConfig(bytes.NewBuffer([]byte(content)))
	err = viper.ReadConfig(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	cfg := config.ServerConfig{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-srv.yaml",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		panic(err)
	}
}
