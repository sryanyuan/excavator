package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

const (
	kBDGUrlSearch = "http://btdigg.pw/"
	kBDGUrlPage   = "http://btdigg.pw/search/"
)

type BDGExecutor struct {
	executeUrl string
	resultSet  []*SearchResult
	TotalPage  int
	Hashkey    string
	TotalCount int
}

func (this *BDGExecutor) Execute(keyword string, maxPage int) error {
	executeUrl := kBDGUrlSearch
	this.resultSet = make([]*SearchResult, 0, 20)

	v := url.Values{
		"keyword": []string{keyword},
	}

	/*body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := new(http.Client)
	req, _ := http.NewRequest("POST", executeUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	log.Println("First request to get data...")
	resp, err := client.Do(req)
	if nil != err {
		log.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Println("Error:", err)
		return err
	}*/
	data, err := httpGet(executeUrl, v)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	if nil == data ||
		len(data) == 0 {
		return errors.New("Request timeout, please retry agin")
	}

	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if nil != err {
		log.Println("Error:", err)
		return errors.New("获取数据超时，请重试")
	}

	//	get pages
	pageLinks := make([]string, 0, 10)
	lastPageText := []rune(doc.Find("div.page-split").Find("span").Text())
	if len(lastPageText) == 0 {
		//return errors.New("获取数据超时，请重试")
		return nil
	}
	lastPageNumberText := string(lastPageText[1 : len(lastPageText)-1])
	lastPageUrlHref, _ := doc.Find("div.page-split").Find("a").First().Attr("href")
	lastPageUrlTpl := []rune(lastPageUrlHref)
	lastPageNumber, err := strconv.Atoi(lastPageNumberText)
	if nil != err {
		return err
	}
	this.TotalPage = lastPageNumber
	//	get hash key
	hashKeys := strings.Split(lastPageUrlHref, "/")
	if len(hashKeys) > 4 {
		this.Hashkey = hashKeys[len(hashKeys)-4]
	}

	//	get total search result count
	searchResultCountText := []rune(doc.Find("a.select").Text())
	if len(searchResultCountText) != 0 {
		searchResultCountNumberText := string(searchResultCountText[3 : len(searchResultCountText)-1])
		this.TotalCount, _ = strconv.Atoi(searchResultCountNumberText)
	}
	if this.TotalCount == 0 {
		return nil
	}

	if 0 != maxPage &&
		lastPageNumber > maxPage {
		lastPageNumber = maxPage
	}

	lastPageUrl := ""
	lastPageSepIndex := 0
	for i := len(lastPageUrlTpl) - 1; i >= 0; i-- {
		if lastPageUrlTpl[i] == rune('/') {
			lastPageSepIndex++
		}

		if lastPageSepIndex >= 3 {
			//	sep
			lastPageUrl = string(lastPageUrlTpl[:i+1])
			break
		}
	}
	if len(lastPageUrl) != 0 {
		for i := 2; i <= lastPageNumber; i++ {
			pageLinks = append(pageLinks, lastPageUrl+strconv.Itoa(i)+"/0/0.html")
		}
	}

	//	parse home page
	this.parseDoc(doc)

	//	parse other pages
	for i, v := range pageLinks {
		log.Println("Get page", i, " data...")
		time.Sleep(time.Millisecond * 100)
		pageClient := new(http.Client)
		pageReq, _ := http.NewRequest("GET", v, nil)
		pageReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		pageReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

		pageResp, err := pageClient.Do(pageReq)
		if nil != err {
			log.Println("Error:", err)
			continue
		}
		defer pageResp.Body.Close()
		pageData, err := ioutil.ReadAll(pageResp.Body)
		if nil != err {
			log.Println("Error:", err)
			continue
		}

		pageDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(pageData))
		if nil != err {
			log.Println("Error:", err)
			continue
		}
		this.parseDoc(pageDoc)
	}

	return nil
}

func (this *BDGExecutor) ExecutePage(keyhash string, page int) bool {
	if nil == this.resultSet {
		this.resultSet = make([]*SearchResult, 0, 20)
	} else {
		this.resultSet = this.resultSet[0:0]
	}
	if len(keyhash) == 0 {
		return false
	}
	if page < 1 {
		return false
	}

	pageUrl := fmt.Sprintf("%s%s/%d/0/0.html", kBDGUrlPage, keyhash, page)
	/*pageClient := new(http.Client)
	pageReq, _ := http.NewRequest("GET", pageUrl, nil)
	pageReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	pageReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	pageResp, err := pageClient.Do(pageReq)
	if nil != err {
		log.Println("Error:", err)
		return false
	}
	defer pageResp.Body.Close()
	pageData, err := ioutil.ReadAll(pageResp.Body)
	if nil != err {
		log.Println("Error:", err)
		return false
	}*/
	pageData, err := httpGet(pageUrl, nil)
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	if nil == pageData ||
		len(pageData) == 0 {
		log.Println("empty body")
		return false
	}

	pageDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(pageData))
	if nil != err {
		log.Println("Error:", err)
		return false
	}
	this.parseDoc(pageDoc)
	return true
}

func (this *BDGExecutor) parseDoc(doc *goquery.Document) {
	nodes := doc.Find(".list").Find("dl")

	//	get search item
	o := otto.New()
	var jsValue otto.Value
	var err error

	nodes.Each(func(i int, sel *goquery.Selection) {
		item := &SearchResult{}
		item.Name = sel.Find("dd.flist").Find("p").Find("span.filename").Text()
		item.CollectTime = sel.Find("span").Eq(0).Find("b").Text()
		item.FileSize = sel.Find("span").Eq(1).Find("b").Text()
		item.FileCount = sel.Find("span").Eq(2).Find("b").Text()
		item.DownloadSpeed = sel.Find("span").Eq(3).Find("b").Text()
		popularStr := sel.Find("span").Eq(4).Find("b").Text()
		item.Popular, _ = strconv.Atoi(popularStr)
		item.MagnetURI = sel.Find("span").Eq(5).Find("script").Text()

		//	process name and magnet
		if strings.Index(item.Name, "document.write(") == 0 {
			strData := []rune(item.Name)
			strData = strData[15 : len(strData)-1-1]
			item.Name = string(strData)
		}
		jsValue, err = o.Run(item.Name)
		if err == nil &&
			jsValue.IsString() {
			item.Name, err = jsValue.ToString()
		}

		if strings.Index(item.MagnetURI, "document.write(") == 0 {
			strData := []rune(item.MagnetURI)
			strData = strData[15 : len(strData)-1-1]
			item.MagnetURI = string(strData)
		}
		jsValue, err = o.Run(item.MagnetURI)
		if err == nil &&
			jsValue.IsString() {
			item.MagnetURI, err = jsValue.ToString()

			begIndex := strings.Index(item.MagnetURI, "href='")
			endIndex := strings.Index(item.MagnetURI, "' >磁力链</a>")

			if begIndex != -1 &&
				endIndex != -1 {
				strData := []rune(item.MagnetURI)
				strData = strData[begIndex+6 : endIndex]
				item.MagnetURI = string(strData)
			}
		}

		//	add to result list
		this.resultSet = append(this.resultSet, item)
	})
}

func (this *BDGExecutor) GetResult() []*SearchResult {
	if nil == this.resultSet {
		return nil
	}

	return this.resultSet
}
