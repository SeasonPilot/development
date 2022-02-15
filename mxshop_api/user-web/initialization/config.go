package initialization

import (
	"fmt"

	"mxshop-api/user-web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
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

	err = v.Unmarshal(global.SrvConfig)
	if err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.String())
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.SrvConfig)
	})
}
