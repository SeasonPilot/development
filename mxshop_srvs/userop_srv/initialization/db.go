package initialization

import (
	"fmt"
	"log"
	"os"
	"time"

	"mxshop-srvs/userop_srv/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	c := global.ServiceConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.Name)
	zap.S().Infof("数据库连接信息: %s", dsn)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		})

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	//err = global.DB.AutoMigrate(&model.Goods{})
	//if err != nil {
	//	panic(err)
	//}
}
