# lixinger
Stock data fetcher base on lixinger's open API https://www.lixinger.com/open/

Examples
```golang
	id := "00700"
	obj := New("YOUR TOKEN from https://www.lixinger.com/open/api/token")

	float, err := obj.GetFinanceMetricsFloat64(id, "2017-12", "q.metrics.roe.ttm")
	if err == nil {
		fmt.Println(float)
	}


	data, err := obj.GetMarketJsonData(id)
	if err == nil {
		fmt.Println(string(data))
	}

	float, err := obj.GetMarketMetricsFloat64(id, "pe_ttm")
	if err == nil {
	  fmt.Println(float)
	}

	str, err := obj.GetMarketMetricsString(id, "stockCnName")
	if err == nil {
		fmt.Println(str)
	}
```
