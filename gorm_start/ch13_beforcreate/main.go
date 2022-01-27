package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Language struct {
	Name    string
	AddTime sql.NullTime
}

// TableName 自定义表名
func (l *Language) TableName() string {
	return "my_language"
}

// BeforeCreate 每个记录创建的时候自动加上当前时间加入到AddTime中
func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {
	l.AddTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	return nil
}

/*
	1. 我们自己定义表名是什么
	2. 统一的给所有的表名加上一个前缀
*/

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
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "byte_", // 为所有创建的表加前缀; 会被 TableName 覆盖; NamingStrategy 和 TableName 不能同时配置
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Language{})
	if err != nil {
		panic(err)
	}

	db.Create(&Language{
		Name: "season",
	})
}
