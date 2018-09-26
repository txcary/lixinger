package lixinger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/txcary/goutils"
	simplejson "github.com/bitly/go-simplejson"
	"strings"
	"sync"
	"time"
)

var (
	financeMetrics = [...]string{
		"q.profitStatement.np.ttm_y2y",
		"q.profitStatement.np.ttm",
		"q.profitStatement.toi.ttm_y2y",
		"q.profitStatement.oc.ttm_y2y",
		"q.balanceSheet.tseattpc.t",
		"q.balanceSheet.tca.t",
		"q.balanceSheet.tcl.t",
		"q.balanceSheet.ta.t_y2y",
		"q.cashFlow.ncffoa.ttm",
		"q.metrics.roe.ttm",
	}
)

type Finance struct {
	financeMap sync.Map
}

func (obj *Lixinger) getFinancePostBody(id string) ([]byte, error) {
	lastYear := time.Now().Year() - 1
	if time.Now().Month() < 3 {
		lastYear--
	}
	postBody := make(map[string]interface{})
	postBody["token"] = obj.token
	postBody["startDate"] = fmt.Sprintf("%d-12-31", lastYear-9)
	postBody["metrics"] = financeMetrics
	postBody["stockCodes"] = []string{id}
	requestBytes, err := json.Marshal(postBody)
	return requestBytes, err
}

func (obj *Lixinger) getFinanceUrl(id string) (string, error) {
	marketType, err := obj.getMarketType(id)
	industryType, err := obj.getIndustryType(id)
	return `https://open.lixinger.com/api/stock/fs/` + marketType + `/` + industryType, err
}

func (obj *Lixinger) getFinanceJsonData(id string) ([]byte, error) {
	var err error
	var v interface{}
	var ok bool
	if v, ok = obj.financeMap.Load(id); !ok {
		url, err := obj.getFinanceUrl(id)
		postBody, err := obj.getFinancePostBody(id)
		data, err := utils.HttpPostJson(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		v,_ = obj.financeMap.LoadOrStore(id, data)
	}

	return v.([]byte), err
}

func (obj *Lixinger) filterFinanceMetricsFloat64(id string, date string, dataMetrics string) ([]float64, error) {
	data, err := obj.getFinanceJsonData(id)
	sjson, err := simplejson.NewJson(data)
	dataArray := sjson.Get(`data`).MustArray()
	metricsArray := strings.Split(dataMetrics, ".")
	res := make([]float64,0)

	if date == "latest" {
		item := sjson.Get(`data`).GetIndex(0)
		for _, metrics := range metricsArray {
			item = item.Get(metrics)
		}
		res = append(res, item.MustFloat64())
		return res, err
	} else {
		for idx, _ := range dataArray {
			item := sjson.Get(`data`).GetIndex(idx)
			itemDate := item.Get(`date`).MustString()
			//if isMatch, _ := regexp.MatchString("^"+date, itemDate); !isMatch {
			if !strings.Contains(itemDate, date) {
				continue
			}
			for _, metrics := range metricsArray {
				item = item.Get(metrics)
			}
			res = append(res, item.MustFloat64())
		}
	}

	return res, err
}

func (obj *Lixinger) getFinanceMetricsFloat64(id string, date string, dataMetrics string) (float64, error) {
	dataArray, err := obj.filterFinanceMetricsFloat64(id, date, dataMetrics)
	if len(dataArray) == 0 {
		return -1, errors.New(date+":"+dataMetrics+" Not found!")
	}
	return dataArray[0], err
}

func (obj *Lixinger) getFinanceMetricsString(id string, date string, dataMetrics string) (string, error) {
	return "", errors.New("Not Supported yet!")
}

func (obj *Lixinger) initFinance() {
	for _, metrics := range financeMetrics {
		obj.strategyMap.Store(metrics, strategyFinance)
	}
}
