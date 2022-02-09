package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	// fixme: 要先 New()
	v := viper.New()

	v.SetConfigFile("debug-config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Get("name"))

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.String())
		_ = v.ReadInConfig()
		fmt.Println(v.Get("name"))
	})

	time.Sleep(time.Second * 300)
}
