package lixinger

import (
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
)

var (
	marketMetrics = [...]string{
		"pe_ttm",
		"pb",
		"dividend_r",
		"stock_price",
	}
)

type Market struct {
	marketMap map[string][]byte
}

func (obj *Lixinger) getMarketUrl(id string) (string, error) {
	marketType, err := obj.getMarketType(id)
	return `https://open.lixinger.com/api/` + marketType + `/stock/fundamental`, err
}

func (obj *Lixinger) getMarketPostBody(id string) ([]byte, error) {
	postBody := make(map[string]interface{})
	postBody["token"] = obj.token
	postBody["date"] = `latest`
	//TODO: Need to get 3 years ProfitDividendRate
	//postBody["startDate"] = fmt.Sprintf("%d-12-31", lastYear-2)
	postBody["metrics"] = marketMetrics
	postBody["stockCodes"] = []string{id}
	requestBytes, err := json.Marshal(postBody)
	return requestBytes, err
}

func (obj *Lixinger) GetMarketJsonData(id string) ([]byte, error) {
	var err error
	if _, ok := obj.marketMap[id]; !ok {
		url, err := obj.getMarketUrl(id)
		fmt.Println(url)
		postBody, err := obj.getMarketPostBody(id)
		data, err := obj.httpPost(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		obj.marketMap[id] = data
	}

	return obj.marketMap[id], err
}

func (obj *Lixinger) GetMarketMetricsFloat64(id string, dataMetrics string) (float64, error) {
	data, err := obj.GetMarketJsonData(id)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustFloat64()
	return metrics, err
}

func (obj *Lixinger) GetMarketMetricsString(id string, dataMetrics string) (string, error) {
	data, err := obj.GetMarketJsonData(id)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustString()
	return metrics, err
}
func (obj *Lixinger) InitMarket() {
	obj.marketMap = make(map[string][]byte)
}
