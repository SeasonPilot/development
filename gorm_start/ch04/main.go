package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

	// 批量插入
	// 单一的 SQL 语句
	users := []User{{Name: "bobby1"}, {Name: "bobby2"}, {Name: "bobby3"}}
	db.Create(&users)

	//为什么不一次性提交所有的 还要分批次，sql语句有长度限制
	db.CreateInBatches(users, 2)

	for _, user := range users {
		fmt.Println(user.ID) // 1,2,3
	}
	// 根据 Map 创建
	db.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "bobby1"},
		{"Name": "bobby2"},
		{"Name": "bobby3"},
	})
}
