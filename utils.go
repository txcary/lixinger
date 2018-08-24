package lixinger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	httpTypeJson string = "application/json"
)

func httpPostJson(postBody []byte, url string) ([]byte, error) {
	var resp *http.Response
	resp, err := http.Post(url, httpTypeJson, bytes.NewReader(postBody))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func stringCount(str string) int {
	return strings.Count(str, "") - 1
}
