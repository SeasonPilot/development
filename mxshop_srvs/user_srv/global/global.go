package global

import (
	"log"
	"os"
	"time"

	"development/mxshop_srvs/user_srv/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	dsn := "golang:golang@2020@tcp(172.19.30.30:3306)/mxshop_user_srv?parseTime=true&loc=Local&charset=utf8mb4"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		})

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}