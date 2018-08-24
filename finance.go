package lixinger

import (
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"regexp"
	"strings"
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
	financeMap map[string][]byte
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
	return `https://open.lixinger.com/api/` + marketType + `/stock/fs/` + industryType, err
}

func (obj *Lixinger) GetFinanceJsonData(id string) ([]byte, error) {
	var err error
	if _, ok := obj.financeMap[id]; !ok {
		url, err := obj.getFinanceUrl(id)
		postBody, err := obj.getFinancePostBody(id)
		data, err := httpPostJson(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		obj.financeMap[id] = data
	}

	return obj.financeMap[id], err
}

func (obj *Lixinger) GetFinanceMetricsFloat64(id string, date string, dataMetrics string) (float64, error) {
	data, err := obj.GetFinanceJsonData(id)
	sjson, err := simplejson.NewJson(data)
	dataArray := sjson.Get(`data`).MustArray()
	metricsArray := strings.Split(dataMetrics, ".")
	for idx, _ := range dataArray {
		item := sjson.Get(`data`).GetIndex(idx)
		itemDate := item.Get(`date`).MustString()
		if isMatch, _ := regexp.MatchString("^"+date, itemDate); !isMatch {
			continue
		}
		for _, metrics := range metricsArray {
			item = item.Get(metrics)
		}
		return item.MustFloat64(), err
	}
	return -1, errors.New("Date not found!")
}

func (obj *Lixinger) initFinance() {
	obj.financeMap = make(map[string][]byte)
}
