package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Code  sql.NullString
	Price uint
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

	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic(err)
	}

	// 新增
	db.Create(&Product{Code: sql.NullString{String: "D42", Valid: true}, Price: 200})

	// Read
	var product Product
	db.First(&product, 2)                 // 根据整形主键查找
	db.First(&product, "Code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// update 前面要有 Model, Model 中的 ID 会作为 WHERE 的条件
	// UPDATE `products` SET `code`='D45',`updated_at`='2022-01-25 15:16:43.499' WHERE `id` = 3 AND `products`.`deleted_at` IS NULL
	db.Model(&Product{Model: gorm.Model{ID: 3}}).Update("Code", "D45")
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Code: sql.NullString{String: "", Valid: true}, Price: 100}) // 通过 NullString 解决不能更新零值的问题

	// Delete - 删除 product， 并没有执行delete语句，逻辑删除
	db.Delete(&Product{}, "ID = ?", 1)
}
