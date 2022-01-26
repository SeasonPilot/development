package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User `User` 属于 `Company`，`CompanyID` 是外键
type User struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company
}

type Company struct {
	ID   int
	Name string
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

	db.Create(&User{
		Name: "season",
		Company: Company{
			Name: "bytes", // 第二次会创建新的公司
		},
	})

	db.Create(&User{
		Name: "season",
		Company: Company{
			ID: 1, // 将 user 插入已有的公司内
		},
	})

	// 官方文档没有看见这个写法，但试了是可以的
	db.Create(&User{
		Name:      "bobby",
		CompanyID: 1,
	})
}
