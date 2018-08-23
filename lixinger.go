package lixinger

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	httpTypeJson string = "application/json"
)

type Lixinger struct {
	Finance
	Market
	StockInfo

	token string
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

func (obj *Lixinger) Init(token string) {
	obj.token = token

	obj.InitFinance()
	obj.InitMarket()
	obj.InitStockInfo()
}

func New(token string) *Lixinger {
	obj := new(Lixinger)
	obj.Init(token)
	return obj
}
