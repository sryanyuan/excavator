package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//"os"
	"strings"

	"time"

	"sort"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

type BDGResult struct {
	Name          string
	DownloadTimes string
	DownloadSpeed string
	CollectTime   string
	FileSize      string
	FileCount     string
	Popular       int
	MagnetURI     string
}

type BDGExecutor struct {
	executeUrl string
	resultSet  []*BDGResult
}

type BDGResultSet []*BDGResult

func (s BDGResultSet) Len() int {
	return len(s)
}

func (s BDGResultSet) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func (s BDGResultSet) Less(i, j int) bool {
	return s[i].Popular < s[j].Popular
}

func (this *BDGExecutor) Execute(executeUrl string, key string, values []string, maxPage int) error {
	this.executeUrl = executeUrl

	v := url.Values{
		key: values,
	}

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
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
	}

	if nil == data ||
		len(data) == 0 {
		return errors.New("Request timeout, please retry agin")
	}
	log.Println("Done...")

	this.resultSet = make([]*BDGResult, 0, 20)

	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if nil != err {
		log.Println("Error:", err)
		return err
	}

	//	get pages
	pageLinks := make([]string, 0, 10)
	lastPageText := []rune(doc.Find("div.page-split").Find("span").Text())
	if len(lastPageText) == 0 {
		return errors.New("Invalid page number")
	}
	lastPageNumberText := string(lastPageText[1 : len(lastPageText)-1])
	lastPageUrlHref, _ := doc.Find("div.page-split").Find("a").First().Attr("href")
	lastPageUrlTpl := []rune(lastPageUrlHref)
	lastPageNumber, err := strconv.Atoi(lastPageNumberText)
	if nil != err {
		return err
	}
	log.Println("Get total page:", lastPageNumber, "page limit:", maxPage)
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

func (this *BDGExecutor) parse(data []byte) error {
	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if nil != err {
		log.Println("Error:", err)
		return err
	}

	this.parseDoc(doc)

	return nil
}

func (this *BDGExecutor) parseDoc(doc *goquery.Document) {
	nodes := doc.Find(".list").Find("dl")

	//	get search item
	o := otto.New()
	var jsValue otto.Value
	var err error

	nodes.Each(func(i int, sel *goquery.Selection) {
		item := &BDGResult{}
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

	//	sort by Popular
	sort.Sort(sort.Reverse(BDGResultSet(this.resultSet)))
}

func (this *BDGExecutor) GetResult() interface{} {
	if nil == this.resultSet {
		return nil
	}

	return this.resultSet
}
