package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	URL    = `http://gd.189.cn/goods/kdts/verifyCertNum.action?goodsCode=10100212&certNum=%s`
	COOKIE = `_gscu_1708861450=38588434421ddf13; _gscbrs_1708861450=1; citrix_ns_id=8YsbRrWX9DfJRKfghhuRIas2bvkA020; citrix_ns_id_.189.cn_%2F_wlf=TlNDX3h1LTIyMi42OC4xODUuMjI5?SFyPH7LPjeXzuvkfet/gspvs/akA&; cityCode=sh; SHOPID_COOKIEID=10003; svid=9C53E807F77D4319; s_fid=0BEAFDA6DEEDCA80-2753121432C05A11; loginStatus=non-logined; lvid=df2a76488be7f960fa574277a1876881; nvid=1; trkId=89AE5033-5B6D-41F9-9870-6C0251AFD384; s_cc=true; LATN_CODE_COOKIE=0755; BIGipServerTongYiSouSuo=3060705472.39455.0000; d_cmpid=Null; d_source=other; d_channel=0; d_appid=null; d_openid=null; i_cc=true; UM_distinctid=1663b06e6535d5-091b6929ce897f-8383268-149c48-1663b06e654f36; i_click=%5B%5BB%5D%5D; TS9d76e8=b3c808d1c9402ca24cf949d682f84c5ca10ba78689c080115bb5217e; JSESSIONID=h9E7itXH8FFDF6NO3hN2_IlZoThsKz3_6QQy8SLoES95LEuSFpXU!397114502; SESSIONID=59321b27-a44a-4c69-b434-943b9a9c27b2; ijg_s=Less%20than%201%20day; i_invisit=1; i_vnum=2; i_PV=gd.189.cn%2Fgoods%2Fhtml%2Fkdts%2Fmain%2FoperIndex.html; i_ppv=71; ecss_identity=41853938569351438125; i_sq=eship-gdt-prd-new%3D%2526pid%253Dgd.189.cn%25252Fgoods%25252Fhtml%25252Fkdts%25252Fmain%25252FoperIndex.html%2526pidt%253D1%2526oid%253Dhttp%25253A%25252F%25252Fgd.189.cn%25252Fgoods%25252Fhtml%25252Fkdts%25252Fmain%25252FoperIndex.html%25253Flink%25253Dwyts%252523%252523%252523%2526ot%253DA; ijg=1538597264983; i_url=http%3A%2F%2Fgd.189.cn%2Fcommon%2FnewLogin%2Findex_right.htm%3FSSOArea%3D0755%26SSOCustType%3D0`
)

func Check(sufix string) {
	date := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local)
	log.Println("约定时间：", date.Year(), int(date.Month()), date.Day(),
		date.Hour(), date.Minute(), date.Second())

	for ; date.Year() == 2018; date = date.Add(time.Hour * 24) {
		dateStr := fillDate(date) + sufix
		url := fmt.Sprintf(URL, dateStr)
		log.Printf("正在验证：%s\n", url)
		body, err := Get(url, "", COOKIE)
		if err != nil {
			log.Printf("网络请求（%s）出错：%s\n", url, err)
			continue
		}

		bs, err := GbkToUtf8([]byte(body))
		if err != nil {
			log.Printf("GBK转UTF8出错：%s\n", err)
			return
		}
		body = string(bs)

		if strings.Contains(body, "false") {
			log.Printf("验证未通过：%s\n", body)
			continue
		}

		log.Printf("验证通过：%s\n", dateStr)
		return
	}
	log.Printf("此年所有日期都已测试：未通过")
}

func fillDate(date time.Time) string {
	var m, d string
	m = FillPrefixWith(int(date.Month()), "0")
	d = FillPrefixWith(date.Day(), "0")
	return m + d
}

func FillPrefixWith(value int, prefix string) string {
	if value >= 10 {
		return strconv.Itoa(value)
	}
	return prefix + strconv.Itoa(value)
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
