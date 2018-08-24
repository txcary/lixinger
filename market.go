package lixinger

import (
	"encoding/json"
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
	date string
}

func (obj *Lixinger) getMarketUrl(id string) (string, error) {
	marketType, err := obj.getMarketType(id)
	return `https://open.lixinger.com/api/` + marketType + `/stock/fundamental`, err
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

func (obj *Lixinger) GetMarketJsonData(id string, date string) ([]byte, error) {
	var err error
	if _, ok := obj.marketMap[id]; !ok || date!=obj.date {
		url, err := obj.getMarketUrl(id)
		postBody, err := obj.getMarketPostBody(id, date)
		data, err := httpPostJson(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		obj.marketMap[id] = data
		obj.date = date
	}
	return obj.marketMap[id], err

}

func (obj *Lixinger) GetMarketMetricsFloat64(id string, date string, dataMetrics string) (float64, error) {
	data, err := obj.GetMarketJsonData(id, date)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustFloat64()
	return metrics, err
}

func (obj *Lixinger) GetMarketMetricsString(id string, date string, dataMetrics string) (string, error) {
	data, err := obj.GetMarketJsonData(id, date)
	sjson, err := simplejson.NewJson(data)
	metrics := sjson.Get(`data`).GetIndex(0).Get(dataMetrics).MustString()
	return metrics, err
}
func (obj *Lixinger) initMarket() {
	obj.marketMap = make(map[string][]byte)
}
