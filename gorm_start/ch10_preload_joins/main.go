package main

import (
	"fmt"
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

	var user User
	db.Preload("Company").First(&user) // 执行两条 SQL 语句
	fmt.Println(user.ID)
	fmt.Println(user.Company)
	fmt.Println(user.CompanyID)
	fmt.Println(user.Company.Name)

	db.Joins("Company").First(&user) // 执行一条 SQL 语句
	fmt.Println(user.ID)
	fmt.Println(user.Company)
	fmt.Println(user.CompanyID)
	fmt.Println(user.Company.Name)
}
