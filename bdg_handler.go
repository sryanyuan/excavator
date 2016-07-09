package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
	//"log"
)

var tplFuncMap = template.FuncMap{
	"writeMagnet": tplfn_writeMagnet,
}

func tplfn_writeMagnet(item *SearchResult) template.URL {
	return template.URL(item.MagnetURI)
}

func pageAjaxHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	keyhash := r.Form.Get("keyhash")
	pageStr := r.Form.Get("page")
	resultJson := []byte("[]")

	defer func() {
		w.Write(resultJson)
	}()

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return
	}

	if len(keyhash) == 0 {
		return
	}

	exe := &BDGExecutor{}
	if !exe.ExecutePage(keyhash, page) {
		return
	}

	result := exe.GetResult()
	resultJson, _ = json.Marshal(result)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	keyword := r.Form.Get("keyword")
	maxpage := r.Form.Get("maxpage")
	maxpageInt, _ := strconv.Atoi(maxpage)
	maxpageInt = 1
	tmNow := time.Now()

	t := template.New("search.tpl").Funcs(tplFuncMap)
	t, err := t.ParseFiles("template/search.tpl")
	if nil != err {
		panic(err)
	}

	data := make(map[string]interface{})
	data["ResultCount"] = 0
	data["Error"] = ""
	data["TotalPage"] = 0
	data["CurrentPage"] = 0
	data["Keyhash"] = ""
	data["ProcessTime"] = 0
	var buf bytes.Buffer

	defer func() {
		data["ProcessTime"] = time.Now().Sub(tmNow).Nanoseconds() / 1e6
		err = t.Execute(&buf, data)
		if nil != err {
			panic(err)
		}
		w.Write(buf.Bytes())
	}()

	if len(keyword) == 0 {
		//	show homepage
		data["LastSearch"] = ""
		data["Error"] = "关键词不能为空"
	} else {
		addSearchRecord(r.RemoteAddr, keyword)

		data["LastSearch"] = keyword

		//	search it
		executor := &BDGExecutor{}
		if err := executor.Execute(keyword, maxpageInt); nil != err {
			data["Error"] = err.Error()
		} else {
			data["TotalPage"] = executor.TotalPage
			data["CurrentPage"] = 1
			data["Keyhash"] = executor.Hashkey
		}
		resultSet := executor.GetResult()
		sort.Sort(sort.Reverse(SearchResultSet(resultSet)))
		data["ResultCount"] = executor.TotalCount
		data["SearchResult"] = resultSet
	}
}
