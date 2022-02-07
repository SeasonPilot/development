package main

import "go.uber.org/zap"

func main() {
	//zap.NewDevelopment()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	// fixme: 退出前要调用 Sync
	defer logger.Sync()

	url := "baidu.com"
	logger.Info("fail to fetch url", zap.Int("nums", 3), zap.String("url", url))

	sugar := logger.Sugar()

	// templated格式  "msg":"fail to fetch url: baidu.com ,nums: 3"
	sugar.Infof("fail to fetch url: %s ,nums: %d", url, 3)

	// 键值对格式   "msg":"fail to fetch url","nums":3,"url":"baidu.com"
	sugar.Infow("fail to fetch url",
		"nums", 3,
		"url", url,
	)
}
