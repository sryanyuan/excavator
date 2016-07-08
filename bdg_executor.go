package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BDGResult struct {
	Name          string
	DownloadTimes string
	CollectTime   string
	FileSize      string
	FileCount     string
	Popular       string
	MagnetURI     string
}

type BDGExecutor struct {
	executeUrl string
	data       []byte
}

func (this *BDGExecutor) Execute(executeUrl string, key string, values []string) error {
	this.executeUrl = executeUrl

	v := url.Values{
		key: values,
	}

	//把form数据编下码
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := new(http.Client)
	req, _ := http.NewRequest("POST", executeUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	resp, err := client.Do(req) //发送
	if nil != err {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Println("Error:", err)
		return err
	}
	this.data = data
	log.Println(string(data))

	return nil
}

func (this *BDGExecutor) Parse() (interface{}, error) {
	reader := bytes.NewBuffer(this.data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if nil != err {
		log.Println("Error:", err)
		return nil, err
	}

	nodes := doc.Find(".attr")
	if nil == nodes {
		log.Println("Error:", err)
		return nil, errors.New("invalid .attr")
	}

	nodes.Each(func(i int, sel *goquery.Selection) {
		log.Println(i, "get", sel.Text())
	})

	return nil, nil
}

func (this *BDGExecutor) GetPageCount() {

}
