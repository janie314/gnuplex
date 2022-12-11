package sqliteconn

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func Init() *sql.DB {
	db, err := sql.Open("sqlite", "../tmp/gnuplex.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table if not exists pos_cache (filepath string not null primary key, pos int);")
	res := true
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	_, err = db.Exec("create table if not exists history (id integer not null unique, mediafile	text, primary key(id AUTOINCREMENT));")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	_, err = db.Exec("create table if not exists medialist (filepath text not null,  primary key(filepath)) ;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	_, err = db.Exec("create table if not exists mediadirs (filepath text not null, primary key(filepath)) ;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	if !res {
		log.Fatal("I could not initialize the database...")
	}
	return db
}

func GetMedialib(db *sql.DB) []string {
	rows, err := db.Query("select filepath from medialist order by filepath;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return []string{}
	}
	// TODO: append or [i]
	res := make([]string, 131072)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			res[i] = str
			i++
		}
	}
	return res[:i]
}

func Last25(db *sql.DB) []string {
	rows, err := db.Query("select distinct mediafile from history order by id desc limit 25;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return []string{}
	}
	res := make([]string, 16384)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			res[i] = str
			i++
		}
	}
	return res[:i]
}
