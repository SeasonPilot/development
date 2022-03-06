package initialization

import (
	"fmt"
	"strings"

	"mxshop-srvs/userop_srv/config"
	"mxshop-srvs/userop_srv/global"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var nacosInfo config.NacosConfig

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// 从本地文件读取 NacosConfig
	InitConfigFromFile()

	// ServerConfig 是 slice
	// 从 nacos 配置中心读取配置
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosInfo.Host,
			Port:   nacosInfo.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         nacosInfo.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
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
		DataId: nacosInfo.DataId,
		Group:  nacosInfo.Group,
	})
	if err != nil {
		panic(err)
	}

	// 从 io.Reader 读取 config
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	err = viper.ReadConfig(strings.NewReader(content))
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&global.ServiceConfig)
	if err != nil {
		panic(err)
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: nacosInfo.DataId,
		Group:  nacosInfo.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)

			content, err = configClient.GetConfig(vo.ConfigParam{
				DataId: nacosInfo.DataId,
				Group:  nacosInfo.Group,
			})
			if err != nil {
				panic(err)
			}

			// 从 io.Reader 读取 config
			viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
			err = viper.ReadConfig(strings.NewReader(content))
			if err != nil {
				panic(err)
			}

			err = viper.Unmarshal(&global.ServiceConfig)
			if err != nil {
				panic(err)
			}

			fmt.Println("nacos 配置更新: ", global.ServiceConfig)
		},
	})
	if err != nil {
		panic(err)
	}
}

func InitConfigFromFile() {
	configFile := "config-pro.yaml"
	if GetEnvInfo("MXSHOP_DEBUG") {
		configFile = "config-debug.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&nacosInfo)
	if err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.String())
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&nacosInfo)
	})
}
