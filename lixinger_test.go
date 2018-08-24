package lixinger

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"runtime"
)

const (
	configFile string = "lixinger.ini"
)

func slash() string {
	var ostype = runtime.GOOS
	if ostype == "windows" {
		return "\\"
	}
	if ostype == "linux" {
		return "/"
	}
	return "/"
}

func ExampleNew() {
	id := "00700"
	gopath := os.Getenv("GOPATH")
	config, err := ini.Load(gopath + slash() + configFile)
	if err != nil {
		panic(err)
	}
	token := config.Section("").Key("token").String()
	obj := New(token)

	name, err := obj.GetMarketMetricsString(id, "latest", "stockCnName")
	if err == nil {
		fmt.Println(name)
	}

	roe, err := obj.GetFinanceMetricsFloat64(id, "2017-12-31", "q.metrics.roe.ttm")
	if err == nil {
		fmt.Println(int(roe * 100))
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
}
