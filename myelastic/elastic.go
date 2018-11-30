package myelastic

import (
	"github.com/olivere/elastic"
	"encoding/json"
	"context"
	"fmt"
	"sync"
	"statistics/card"
)

type EsServerInfo struct {
	Ip      string
	Port    string
	Client  *elastic.Client
	EsIndex string
	EsType  string
	Wg		*sync.WaitGroup
}

func NewEsServerInfo (ip ,port,esIndex ,esType string,
	wg *sync.WaitGroup) *EsServerInfo{
	return &EsServerInfo{
		Ip:      ip,
		Port:    port,
		EsIndex: esIndex,
		EsType:  esType,
		Wg:      wg,
	}
}

func (this *EsServerInfo) Init() {
	//ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", this.Ip, this.Port)),
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	this.Client = client
}

func (this *EsServerInfo) Insert(row *card.EsCard) bool {

	data, err := json.Marshal(row)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("保存失败：", data)
		return false
	}
	//fmt.Println("insert data:",string(data))
	exists, err := this.Client.IndexExists(this.EsIndex).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err := this.Client.CreateIndex(this.EsIndex).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	_, err = this.Client.Index().
		Index(this.EsIndex).
		Type(this.EsType).
		Id(row.MakceId()).
		BodyString(string(data)).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	this.Wg.Done()
	return true
}

func (this *EsServerInfo) Update(row *card.EsCard) bool {

	data, err := json.Marshal(row)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("保存失败：", data)
		return false
	}
	//fmt.Println("update data:",string(data))
	_, err = this.Client.Index().
		Index(this.EsIndex).
		Type(this.EsType).
		Id(row.MakceId()).
		BodyString(string(data)).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	this.Wg.Done()
	return true
}

func (this *EsServerInfo) QueryById(row *card.EsCard) bool{
	//字段相等
	_, err := this.Client.Get().Index(this.EsIndex).Type(this.EsType).Id(row.MakceId()).Do(context.Background())
	if err != nil {
		return false
	}
	return true
}

//func (this* EsServerInfo)QueryByDate(date projectType.QueryDate,rowinfos*[]projectType.RowInfo){
//	if !date.IsRange {
//		q := elastic.NewQueryStringQuery(fmt.Sprintf("case_date:%s",date.Begin))
//		res, err := this.Client.Search(this.EsIndex).Type(this.EsType).Size(this.EsRetSize).Query(q).Do(context.Background())
//		if err != nil {
//			println(err.Error())
//		}
//		if res.Hits.TotalHits > 0 {
//			for _, hit := range res.Hits.Hits {
//				var row projectType.RowInfo
//				//fmt.Println(*hit.Source)
//				err := json.Unmarshal(*hit.Source, &row) //另外一种取数据的方法
//				if err != nil {
//					fmt.Println("Deserialization failed")
//				}
//				*rowinfos = append(*rowinfos,row)
//			}
//		} else {
//			fmt.Printf("no found %s case\n",date.Begin)
//		}
//	}else{
//		//boolQ.Must(elastic.NewMatchQuery("last_name", "smith")
//		boolQ := elastic.NewRangeQuery("case_date").Gte(date.Begin).Lte(date.End)
//		res, err := this.Client.Search(this.EsIndex).Type(this.EsIndex).Size(this.EsRetSize).Query(boolQ).Do(context.Background())
//		if err == nil{
//			if res.Hits.TotalHits > 0 {
//				fmt.Printf("Found a total of %d case \n", res.Hits.TotalHits)
//				for _, hit := range res.Hits.Hits {
//					var row projectType.RowInfo
//					err := json.Unmarshal(*hit.Source, &row) //另外一种取数据的方法
//					if err != nil {
//						fmt.Println("Deserialization failed")
//					}
//					*rowinfos = append(*rowinfos,row)
//				}
//			} else {
//				fmt.Println("res.Hits.TotalHits:",res.Hits.TotalHits)
//				fmt.Printf("no found %s~%s case\n",date.Begin,date.End)
//			}
//		}
//	}
//
//}