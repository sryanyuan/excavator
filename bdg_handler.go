package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"log"

	"github.com/golang/protobuf/proto"
	"github.com/sryanyuan/excavator/protocol"
)

var tplFuncMap = template.FuncMap{
	"writeMagnet": tplfn_writeMagnet,
}

func tplfn_writeMagnet(item *BDGResult) template.URL {
	return template.URL(item.MagnetURI)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	keyword := r.Form.Get("keyword")
	maxpage := r.Form.Get("maxpage")
	maxpageInt, _ := strconv.Atoi(maxpage)

	t := template.New("search.tpl").Funcs(tplFuncMap)
	t, err := t.ParseFiles("template/search.tpl")
	if nil != err {
		panic(err)
	}

	data := make(map[string]interface{})
	data["ResultCount"] = 0
	var buf bytes.Buffer

	defer func() {
		err = t.Execute(&buf, data)
		if nil != err {
			panic(err)
		}
		w.Write(buf.Bytes())
	}()

	if len(keyword) == 0 {
		//	show homepage
		data["LastSearch"] = ""
	} else {
		data["LastSearch"] = keyword

		//	search it
		executor := &BDGExecutor{}
		if err := executor.Execute("http://btdigg.pw/", "keyword", []string{keyword}, maxpageInt); nil != err {
			data["Error"] = err.Error()
			return
		}
		result := executor.GetResult()
		resultSet, ok := result.([]*BDGResult)
		if !ok {
			panic("invalid result set")
		}
		if nil == result {
			data["Error"] = "没有搜索结果"
			return
		} else {
			data["ResultCount"] = len(resultSet)
			data["SearchResult"] = resultSet
		}
	}
}

func bdgHandler_(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	resultProto := &protocol.SearchResult{}
	resultProto.Result = proto.Int(1)

	eidStr := r.Form.Get("eid")
	eid, _ := strconv.Atoi(eidStr)
	keyword := r.Form.Get("keyword")
	maxpage := r.Form.Get("maxpage")
	maxpageInt, _ := strconv.Atoi(maxpage)
	format := r.Form.Get("format")

	defer func() {
		var data []byte
		if format == "json" {
			w.Header().Set("Content-Type", "application/json")
			data, _ = json.Marshal(resultProto)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
			data, _ = proto.Marshal(resultProto)
		}

		w.Write(data)
	}()

	if len(keyword) == 0 {
		resultProto.Result = proto.Int(2)
		return
	}

	if 0 == eid {
		executor := &BDGExecutor{}
		if err := executor.Execute("http://btdigg.pw/", "keyword", []string{keyword}, maxpageInt); nil != err {
			log.Println("execute failed, error:", err)
			return
		}
		result := executor.GetResult()
		if nil == result {
			log.Println("empty result")
			return
		}

		if resultSet, ok := result.([]*BDGResult); ok {
			resultProto.Result = proto.Int(0)
			resultProto.ResultSet = make([]*protocol.SearchItem, 0, len(resultSet))

			for _, v := range resultSet {
				ritem := &protocol.SearchItem{}
				ritem.Name = proto.String(v.Name)
				ritem.CollectTime = proto.String(v.CollectTime)
				ritem.DownloadSpeed = proto.String(v.DownloadSpeed)
				downloadTimes, _ := strconv.Atoi(v.DownloadTimes)
				ritem.DownloadTimes = proto.Int(downloadTimes)
				ritem.FileCount = proto.String(v.FileCount)
				ritem.FileSize = proto.String(v.FileSize)
				ritem.MagnetURI = proto.String(v.MagnetURI)
				ritem.Popular = proto.Int(v.Popular)
				resultProto.ResultSet = append(resultProto.ResultSet, ritem)
			}
		}
	}
}
