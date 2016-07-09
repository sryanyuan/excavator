package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SearchResult struct {
	Name          string
	DownloadTimes string
	DownloadSpeed string
	CollectTime   string
	FileSize      string
	FileCount     string
	Popular       int
	MagnetURI     string
}

type SearchResultSet []*SearchResult

func (s SearchResultSet) Len() int {
	return len(s)
}

func (s SearchResultSet) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func (s SearchResultSet) Less(i, j int) bool {
	return s[i].Popular < s[j].Popular
}

type IExecutor interface {
	Execute(key string, maxpage int) error
	GetResult() []*SearchResult
}

func httpGet(urlStr string, values url.Values) ([]byte, error) {
	var body io.ReadCloser
	if nil != values {
		body = ioutil.NopCloser(strings.NewReader(values.Encode()))
	}
	client := new(http.Client)
	req, _ := http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return data, nil
}
