package main

import (
	"flag"
	"log"

	"net/http"
)

var (
	g_recordHistory = false
)

func main() {
	lisaddr := flag.String("lisaddr", "localhost:1111", "listen address")
	recordHistory := flag.Int("history", 0, "record history")
	flag.Parse()

	if *recordHistory != 0 {
		g_recordHistory = true
		initDB()
	}

	//	static file
	http.Handle("/static/css/", http.FileServer(http.Dir(".")))
	http.Handle("/static/js/", http.FileServer(http.Dir(".")))
	http.Handle("/static/img/", http.FileServer(http.Dir(".")))
	http.Handle("/static/fonts/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", searchHandler)
	http.HandleFunc("/page", pageAjaxHandler)
	log.Println("visit", *lisaddr, "to search magnet link")
	http.ListenAndServe(*lisaddr, nil)
}
