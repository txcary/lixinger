package lixinger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	//"fmt"
)

const (
	httpTypeJson string = "application/json"
)

func httpPostJson(postBody []byte, url string) ([]byte, error) {
	//fmt.Println(string(postBody))
	//fmt.Println(url)
	var resp *http.Response
	resp, err := http.Post(url, httpTypeJson, bytes.NewReader(postBody))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	return data, err
}

func stringCount(str string) int {
	return strings.Count(str, "") - 1
}
