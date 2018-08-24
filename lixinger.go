package lixinger

type Lixinger struct {
	Finance
	Market
	StockInfo

	token string
}

func (obj *Lixinger) Init(token string) {
	obj.token = token

	obj.initFinance()
	obj.initMarket()
	obj.initStockInfo()
}

func New(token string) *Lixinger {
	obj := new(Lixinger)
	obj.Init(token)
	return obj
}
