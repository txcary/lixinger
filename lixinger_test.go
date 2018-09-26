package lixinger

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/txcary/goutils"
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
	utils.SetLogLevel(utils.LogLevelDebug)	
	id := "00700"
	gopath := os.Getenv("GOPATH")
	config, err := ini.Load(gopath + slash() + configFile)
	if err != nil {
		panic(err)
	}
	token := config.Section("").Key("token").String()
	obj := New(token)

	name, err := obj.GetString(id, "latest", "stockCnName")
	if err == nil {
		fmt.Println(name)
	}

	//roe, err := obj.GetFloat64(id, "latest", "q.metrics.roe.ttm")
	roe, err := obj.GetFloat64(id, "2017-12-31", "q.metrics.roe.ttm")
	if err == nil {
		fmt.Println(int(roe * 100))
	}

	roeArray, err := obj.FilterFloat64(id, "12-31", "q.metrics.roe.ttm")
	if err == nil {
		for idx,_ := range roeArray {
			fmt.Println(int(roeArray[idx] * 100))
		}
	}

	//pe, err := obj.GetFloat64(id, "latest", "pe_ttm")
	pe, err := obj.GetFloat64(id, "2017-12-29", "pe_ttm")
	if err == nil {
		fmt.Println(int(pe))
	}
	
	//output:
	//腾讯控股
	//31
	//31
	//26
	//28
	//33
	//30
	//35
	//40
	//47
	//53
	//45
	//55
}
