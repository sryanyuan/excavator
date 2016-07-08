package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/sryanyuan/excavator/protocol"
)

func main() {
	flag.Parse()
	lisaddr := flag.String("lisaddr", "localhost:1111", "listen address")

	http.HandleFunc("/s", bdgHandler)
	http.ListenAndServe(*lisaddr, nil)
}

func cmdMain() {
	eId := flag.Int("eid", 0, "executor id")
	keyword := flag.String("keyword", "", "keyword")
	maxpage := flag.Int("maxpage", 0, "max page limit")
	fileId := flag.Int("fileid", 0, "file id")

	flag.Parse()

	resultProto := &protocol.SearchResult{}
	resultProto.Result = proto.Int(1)

	defer func() {
		file := *fileId
		if 0 == file {
			return
		}

		os.Mkdir("search", os.ModeDir)
		f, err := os.Create("search/" + strconv.Itoa(file))
		if nil != err {
			log.Println(err)
			return
		}
		defer f.Close()

		data, _ := proto.Marshal(resultProto)
		f.Write(data)
	}()

	if len(*keyword) == 0 {
		log.Println("Invalid keyword")
		return
	}

	if 0 == *eId {
		executor := &BDGExecutor{}
		if err := executor.Execute("http://btdigg.pw/", "keyword", []string{*keyword}, *maxpage); nil != err {
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

			for i, v := range resultSet {
				log.Println(i, v)
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
