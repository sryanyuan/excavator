package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	g_dbPool *sql.DB
)

func initDB() {
	if !g_recordHistory {
		return
	}

	var err error
	os.Mkdir("history", os.ModeDir)
	g_dbPool, err = sql.Open("sqlite3", "history/search.db")
	if nil != err {
		log.Println("Can't create db")
		return
	}

	expr := `
		CREATE TABLE IF NOT EXISTS user_search_record(id integer primary key autoincrement,
		keyword varchar(255),
		search_time integer,
		ip varchar(15))
	`
	_, err = g_dbPool.Exec(expr)
	if nil != err {
		return
	}
}

func addSearchRecord(ip string, keyword string) {
	if !g_recordHistory {
		return
	}
	g_dbPool.Exec("INSERT INTO user_search_record (keyword, search_time, ip) VALUES (?,?,?)", keyword, time.Now().Unix(), ip)
}
