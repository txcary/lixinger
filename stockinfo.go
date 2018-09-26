package lixinger

import (
	"encoding/json"
	"errors"
	"sync"
	"fmt"
	"github.com/txcary/goutils"
	simplejson "github.com/bitly/go-simplejson"
)

const (
	stockInfoUrl string = `https://open.lixinger.com/api/stock`
)

var (
	industryTypeArray = [...]string{"bank", "insurance", "security"}
	industryAreaArray = [...]string{"cn", "hk"}
)

type StockInfo struct {
	industryMap sync.Map
}

func (obj *Lixinger) initIndustryMap() {
	postBody := make(map[string]interface{})
	postBody["token"] = obj.token

	for _, industryArea := range industryAreaArray {
		postBody["areaCode"] = industryArea
		for _, industryType := range industryTypeArray {
			postBody["industryType"] = industryType
			requestBytes, err := json.Marshal(postBody)
			data, err := utils.HttpPostJson(requestBytes, stockInfoUrl)
			if err == nil {
				sjson, err := simplejson.NewJson(data)
				if err != nil {
					panic(err)
				}
				dataArray := sjson.Get(`data`).MustArray()
				for idx, _ := range dataArray {
					id := sjson.Get(`data`).GetIndex(idx).Get(`stockCode`).MustString()
					obj.industryMap.Store(id, industryType)
				}
			}
		}
	}
}

func (obj *Lixinger) getMarketType(id string) (string, error) {
	marketType := ""
	if utils.StringCount(id) == 5 {
		marketType = "h"
	} else if utils.StringCount(id) == 6 {
		marketType = "a"
	} else {
		errmsg := "Error: ID not correct!"
		fmt.Println(errmsg)
		return marketType, errors.New(errmsg)
	}
	return marketType, nil
}

func (obj *Lixinger) getIndustryType(id string) (string, error) {
	industryType := "industry"
	if newType, ok := obj.industryMap.Load(id); ok {
		industryType = newType.(string)
	}
	return industryType, nil
}

func (obj *Lixinger) initStockInfo() {
	obj.initIndustryMap()
}
