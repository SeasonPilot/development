package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 有多张 CreditCard，UserID 是外键
type User struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	dsn := "golang:golang@2020@tcp(172.19.30.30:3306)/gorm_test?parseTime=true&loc=Local&charset=utf8mb4"

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

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
