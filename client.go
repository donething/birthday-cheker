package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略https证书错误
}
var client = &http.Client{
	Transport: tr,
	Timeout:   50 * time.Second, // 网络超时时间
}

// 执行get请求
func Get(url string, agent string, cookie string) (bodyString string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	if strings.Trim(agent, " ") != "" {
		req.Header.Set("User-Agent", agent)
	}
	if strings.Trim(cookie, " ") != "" {
		req.Header.Set("Cookie", cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	bodyString = string(body)
	return
}
