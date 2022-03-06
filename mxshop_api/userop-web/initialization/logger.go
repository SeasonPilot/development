package initialization

import "go.uber.org/zap"

// InitLogger 设置全局 logger
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 全局的 logger 不需要 Sync ?
	zap.ReplaceGlobals(logger)
}
