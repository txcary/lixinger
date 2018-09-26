package lixinger

import (
	"encoding/json"
	"sync"
	"github.com/txcary/goutils"
	simplejson "github.com/bitly/go-simplejson"
)

var (
	marketMetrics = [...]string{
		"pe_ttm",
		"pb",
		"dividend_r",
		"stock_price",
	}
	marketImplicitMetrics = [...]string{
		"stockCnName",
	}

)

type Market struct {
	marketMap sync.Map
	date      string
}

func (obj *Lixinger) getMarketUrl(id string) (string, error) {
	marketType, err := obj.getMarketType(id)
	return `https://open.lixinger.com/api/stock/fundamental/` + marketType, err
}

func (obj *Lixinger) getMarketPostBody(id string, date string) ([]byte, error) {
	postBody := make(map[string]interface{})
	postBody["token"] = obj.token
	postBody["date"] = date
	postBody["metrics"] = marketMetrics
	postBody["stockCodes"] = []string{id}
	requestBytes, err := json.Marshal(postBody)
	return requestBytes, err
}

func (obj *Lixinger) getMarketJsonData(id string, date string) ([]byte, error) {
	var err error
	var v interface{}
	var ok bool
	if v, ok = obj.marketMap.Load(id); !ok || date != obj.date {
		url, err := obj.getMarketUrl(id)
		postBody, err := obj.getMarketPostBody(id, date)
		data, err := utils.HttpPostJson(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		obj.marketMap.Store(id, data)
		v,_ = obj.marketMap.Load(id)
		obj.date = date
	}
	return v.([]byte), err

}

func (obj *Lixinger) getMarketMetricsFloat64(id string, date string, dataMetrics string) (float64, error) {
	data, err := obj.getMarketJsonData(id, date)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustFloat64()
	return metrics, err
}

func (obj *Lixinger) getMarketMetricsString(id string, date string, dataMetrics string) (string, error) {
	data, err := obj.getMarketJsonData(id, date)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustString()
	return metrics, err
}
func (obj *Lixinger) initMarket() {
	for _, metrics := range marketMetrics {
		obj.strategyMap.Store(metrics, strategyMarket)
	}
	for _, metrics := range marketImplicitMetrics {
		obj.strategyMap.Store(metrics, strategyMarket)
	}
}
