package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

func GetMedialib(db *sql.DB) []string {
	rows, err := db.Query(`select filepath from medialist`)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getmedialib query problem", err)
		return []string{}
	}
	res := make([]string, 16384)
	var line string
	for rows.Next() {
		err = rows.Scan(&line)
		if err != nil {
			res = append(res, line)
		}
	}
	return res
}

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := sql.Open("sqlite", "../tmp/test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	res := initDB(db)
	if !res {
		log.Fatal("I could not initialize the database...")
	}
	for err != nil {
		err = db.Ping()
	}
}

func initDB(db *sql.DB) bool {
	_, err := db.Query("create table if not exists pos_cache (filepath string not null primary key, pos int);")
	res := true
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	db.Query("create table if not exists history (id integer not null unique, mediafile	text, primary key(id AUTOINCREMENT));")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	db.Query("create table if not exists medialist (filepath text not null,  primary key(filepath)) ;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	return res
}
