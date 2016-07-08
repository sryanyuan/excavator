package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"log"

	"github.com/golang/protobuf/proto"
	"github.com/sryanyuan/excavator/protocol"
)

func bdgHandler(w http.ResponseWriter, r *http.Request) {
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
				popular, _ := strconv.Atoi(v.Popular)
				ritem.Popular = proto.Int(popular)
				resultProto.ResultSet = append(resultProto.ResultSet, ritem)
			}
		}
	}
}
