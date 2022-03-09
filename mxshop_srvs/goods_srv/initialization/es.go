package initialization

import (
	"context"
	"fmt"
	"log"
	"os"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"

	"github.com/olivere/elastic/v7"
)

func InitES() {
	url := fmt.Sprintf("http://%s:%d", global.ServiceConfig.EsInfo.Host, global.ServiceConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "mx", log.LstdFlags)

	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	exists, err := global.EsClient.IndexExists(model.EsGoods{}.Name).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := global.EsClient.CreateIndex(model.EsGoods{}.Name).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			fmt.Println(err)
			panic(err)
		}
	}
}
