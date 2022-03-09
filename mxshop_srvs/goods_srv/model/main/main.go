package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"mxshop-srvs/goods_srv/model"

	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	//dsn := "golang:golang@2020@tcp(172.19.30.30:3306)/goods_test?parseTime=true&loc=Local&charset=utf8mb4"
	//
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold:             time.Second,
	//		Colorful:                  true,
	//		IgnoreRecordNotFoundError: false,
	//		LogLevel:                  logger.Info,
	//	})
	//
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = db.AutoMigrate(&model.Category{},
	//	&model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})
	//if err != nil {
	//	panic(err)
	//}
	Mysql2ES()
}

func Mysql2ES() {
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

	url := "http://172.19.30.30:9200/"
	l := log.New(os.Stdout, "mx", log.LstdFlags)

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false), elastic.SetTraceLog(l))
	if err != nil {
		panic(err)
	}

	var goods []model.Goods
	db.Find(&goods)

	for _, g := range goods {
		esGoods := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandsID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}

		_, err = client.Index().Index(esGoods.GetIndexName()).BodyJson(esGoods).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
