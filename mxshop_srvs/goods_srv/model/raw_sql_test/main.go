package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"mxshop-srvs/goods_srv/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "golang:golang@2020@tcp(172.19.30.30:3306)/mxshop_goods_srv?parseTime=true&loc=Local&charset=utf8mb4"

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

	// SELECT * FROM `goods` WHERE category_id in ( select id from category where parent_category_id in (select id from category where parent_category_id = 130358) )
	db = db.Where(fmt.Sprintf(
		"category_id in ( select id from category where parent_category_id in (select id from category where parent_category_id = %d) )",
		130358))
	var goods []model.Goods
	_ = db.Preload("Category").Preload("Brands").Select("id").Find(&goods)

	//
	//q := fmt.Sprintf(" select id from category where parent_category_id in (select id from category where parent_category_id = %d) ", 130358)
	q2 := fmt.Sprintf("select id from category where parent_category_id = %d ", 136781)
	//q3 := fmt.Sprintf("select id from category where id = %d", 238007)

	type Result struct {
		ID int32
	}
	var results []Result
	db.Model(&model.Category{}).Raw(q2).Scan(&results)
	fmt.Println(results)
}
