package lixinger

import (
	"errors"
)

const (
	strategyFinance = iota
	strategyMarket
)

type Lixinger struct {
	Finance
	Market
	StockInfo

	token string
	strategyMap map[string]int
}

func (obj *Lixinger) Init(token string) {
	obj.token = token
	obj.strategyMap = make(map[string]int)

	obj.initFinance()
	obj.initMarket()
	obj.initStockInfo()
}

func (obj *Lixinger) FilterFloat64(id string, date string, dataMetrics string) (res []float64, err error) {
	if strategy,ok := obj.strategyMap[dataMetrics]; ok {
		if strategy == strategyFinance {
			res, err = obj.filterFinanceMetricsFloat64(id, date, dataMetrics)
			return
		}
	}
	err = errors.New(dataMetrics+" not recognized!")
	return
}

func (obj *Lixinger) GetFloat64(id string, date string, dataMetrics string) (res float64, err error) {
	if strategy,ok := obj.strategyMap[dataMetrics]; ok {
		if strategy == strategyFinance {
			res, err = obj.getFinanceMetricsFloat64(id, date, dataMetrics)
			return
		}
		if strategy == strategyMarket {
			res, err = obj.getMarketMetricsFloat64(id, date, dataMetrics)
			return
		}
	}
	err = errors.New(dataMetrics+" not recognized!")
	return
}

func (obj *Lixinger) GetString(id string, date string, dataMetrics string) (res string, err error) {
	if strategy,ok := obj.strategyMap[dataMetrics]; ok {
		if strategy == strategyFinance {
			res, err = obj.getFinanceMetricsString(id, date, dataMetrics)
			return
		}
		if strategy == strategyMarket {
			res, err = obj.getMarketMetricsString(id, date, dataMetrics)
			return
		}
	}
	err = errors.New(dataMetrics+" not recognized!")
	return
}

func New(token string) *Lixinger {
	obj := new(Lixinger)
	obj.Init(token)
	return obj
}
