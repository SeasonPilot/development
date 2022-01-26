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

// User3 拥有并属于多种 language，`user_languages` 是连接表
type User3 struct {
	gorm.Model
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
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

	err = db.AutoMigrate(&User3{})
	if err != nil {
		panic(err)
	}

	db.Create(&User3{
		Languages: []Language{
			{
				Name: "Python",
			},
			{
				Name: "Rust",
			},
		},
	})

	var user = User3{
		Model: gorm.Model{
			ID: 2,
		},
	}
	// 执行 3 条 SQL 语句
	db.Preload("Languages").Find(&user)
	for _, language := range user.Languages {
		fmt.Println(language.Name)
	}

	// 执行 1 条 SQL 语句
	var languages []Language
	err = db.Model(&user).Association("Languages").Find(&languages)
	if err != nil {
		panic(err)
	}
	for _, language := range languages { // 结果写到 languages
		fmt.Println(language.Name)
	}
}
