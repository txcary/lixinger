# lixinger
Stock data fetcher base on lixinger's open API https://www.lixinger.com/open/

Examples
```golang
	id := "00700"
	token := ""YOUR TOKEN from https://www.lixinger.com/open/api/token"
	obj := New(token)

	name, err := obj.GetMarketMetricsString(id, "latest", "stockCnName")
	if err == nil {
		fmt.Println(name)
	}

	roe, err := obj.GetFinanceMetricsFloat64(id, "2017-12-31", "q.metrics.roe.ttm")
	if err == nil {
		fmt.Println(int(roe*100))
	}

	/*
	data, err := obj.GetMarketJsonData(id, "2017-12-29")
	if err == nil {
		fmt.Println(string(data))
	}
	*/

	//pe, err := obj.GetMarketMetricsFloat64(id, "latest", "pe_ttm")
	pe, err := obj.GetMarketMetricsFloat64(id, "2017-12-29", "pe_ttm")
	if err == nil {
	  fmt.Println(int(pe))
	}

	//output:
	//腾讯控股
	//31
	//55
```
