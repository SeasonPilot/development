package main

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./mylog.log",
		"stderr",
		"stdout",
	}
	return cfg.Build()
}

func main() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	su := logger.Sugar()
	defer func(su *zap.SugaredLogger) {
		err = su.Sync()
		if err != nil {
			panic(err)
		}
	}(su)

	url := "baidu.com"
	su.Infof("fail to fetch url %s,%d", url, 3)
}
