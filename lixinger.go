package lixinger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	httpTypeJson string = "application/json"
	stockInfoUrl string = `https://open.lixinger.com/api/stock`
)

var (
	industryTypeArray = [...]string{"bank", "insurance", "security"}
	industryAreaArray = [...]string{"cn", "hk"}
	financeMetrics    = [...]string{
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

type Lixinger struct {
	Id         string
	token      string
	financeMap map[string][]byte
	marketMap  map[string][]byte

	industryMap map[string]string
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

func (obj *Lixinger) httpPost(postBody []byte, url string) ([]byte, error) {
	var resp *http.Response
	resp, err := http.Post(url, httpTypeJson, bytes.NewReader(postBody))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func (obj *Lixinger) initIndustryMap() {
	postBody := make(map[string]interface{})
	postBody["token"] = obj.token

	for _, industryArea := range industryAreaArray {
		postBody["areaCode"] = industryArea
		for _, industryType := range industryTypeArray {
			postBody["industryType"] = industryType
			requestBytes, err := json.Marshal(postBody)
			data, err := obj.httpPost(requestBytes, stockInfoUrl)
			if err == nil {
				sjson, err := simplejson.NewJson(data)
				if err != nil {
					panic(err)
				}
				dataArray := sjson.Get(`data`).MustArray()
				for idx, _ := range dataArray {
					id := sjson.Get(`data`).GetIndex(idx).Get(`stockCode`).MustString()
					obj.industryMap[id] = industryType
				}
			}
		}
	}
}

func (obj *Lixinger) stringCount(str string) int {
	return strings.Count(str, "") - 1
}

func (obj *Lixinger) getFinanceUrl(id string) (string, error) {
	marketType := "a"
	if obj.stringCount(id) == 5 {
		marketType = "h"
	}

	industryType := "industry"
	if newType, ok := obj.industryMap[id]; ok {
		industryType = newType
	}

	return `https://open.lixinger.com/api/` + marketType + `/stock/fs/` + industryType, nil
}

func (obj *Lixinger) GetFinanceJsonData(id string) ([]byte, error) {
	var err error
	if _, ok := obj.financeMap[id]; !ok {
		url, err := obj.getFinanceUrl(id)
		fmt.Println(url)
		postBody, err := obj.getFinancePostBody(id)
		data, err := obj.httpPost(postBody, url)
		if err != nil {
			return []byte{}, err
		}
		obj.financeMap[id] = data
	}

	return obj.financeMap[id], err
}

func (obj *Lixinger) GetFinanceMetrics(id string, date string, dataMetrics string) (float64, error) {
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

func New(token string) *Lixinger {
	obj := new(Lixinger)
	obj.token = token
	obj.financeMap = make(map[string][]byte)
	obj.marketMap = make(map[string][]byte)
	obj.industryMap = make(map[string]string)
	obj.initIndustryMap()

	return obj
}
