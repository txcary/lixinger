# lixinger
Stock data fetcher base on lixinger's open API https://www.lixinger.com/open/

Examples
```golang
	id := "00700"
	token := ""YOUR TOKEN from https://www.lixinger.com/open/api/token"
	obj := New(token)

	name, err := obj.GetString(id, "latest", "stockCnName")
	if err == nil {
		fmt.Println(name)
	}

	roe, err := obj.GetFloat64(id, "2017-12-31", "q.metrics.roe.ttm")
	if err == nil {
		fmt.Println(int(roe * 100))
	}

	//pe, err := obj.GetFloat64(id, "latest", "pe_ttm")
	pe, err := obj.GetFloat64(id, "2017-12-29", "pe_ttm")
	if err == nil {
		fmt.Println(int(pe))
	}


	//output:
	//腾讯控股
	//31
	//55
```
