package main

import (
	"log"
	"os"
	"time"

	"mxshop-srvs/order_srv/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "golang:golang@2020@tcp(172.19.30.30:3306)/mxshop_order_srv?parseTime=true&loc=Local&charset=utf8mb4"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.ShoppingCart{}, &model.OrderInfo{}, &model.OrderGoods{})
	if err != nil {
		panic(err)
	}
}
