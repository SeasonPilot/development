package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

type Account struct {
	AccountNumber int32  `json:"account_number"`
	FirstName     string `json:"firstname"`
}

func main() {
	url := "http://172.19.30.30:9200/"
	l := log.New(os.Stdout, "mx", log.LstdFlags)

	// 这里必须将sniff设置为false，因为使用olivere/elastic连接elasticsearch时，发现连接地址明明输入的时候是外网地址
	// 但是连接时会自动转换成内网地址或者docker中的ip地址，导致服务连接不上。
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false), elastic.SetTraceLog(l))
	if err != nil {
		panic(err)
	}

	// match 查询
	q := elastic.NewMatchQuery("address", "street")
	result, err := client.Search().Index("account").Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Hits.TotalHits.Value)

	for _, hit := range result.Hits.Hits {
		var account Account
		err = json.Unmarshal(hit.Source, &account)
		if err != nil {
			panic(err)
		}
		fmt.Println(account)
	}
}
