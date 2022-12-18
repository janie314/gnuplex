package sqliteconn

import (
	"database/sql"
	"fmt"
	"gnuplex-backend/consts"
	"log"
	"os"
	"strconv"
	"sync"

	_ "modernc.org/sqlite"
)

func Init(mu *sync.Mutex) *sql.DB {
	fmt.Println("got h lock")
	defer fmt.Println("rem h lock")
	mu.Lock()
	defer mu.Unlock()
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
	_, err = db.Exec("create table if not exists version_info (key string not null primary key, value string);")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	_, err = db.Exec("insert or ignore into version_info values ('db_schema_version', ?);", consts.DBVersion)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res = false
	}
	if !res {
		log.Fatal("I could not initialize the database...")
	}
	UpgradeDB(db)
	return db
}

func UpgradeDB(db *sql.DB) error {
	rows, err := db.Query("select value from version_info where key = 'db_schema_version';")
	if err != nil {
		log.Fatal("Upgrade db error", err)
	}
	next := rows.Next()
	if !next {
		log.Fatal("Upgrade db error: no version schema", err)
	}
	var vers string
	err = rows.Scan(&vers)
	if err != nil {
		log.Fatal("Bad version schema 1", err)
	}
	versNum, err := strconv.Atoi(vers)
	if err != nil {
		log.Fatal("Bad version schema 2", err)
	}
	if versNum < consts.DBVersion {
		fmt.Println("Should I do something?")
	}
	rows.Close()
	return nil
}
