package liteDB

import (
	"database/sql"
	"fmt"
	"gnuplex-backend/consts"
	"log"
	"strconv"
	"sync"

	_ "modernc.org/sqlite"
)

type LiteDB struct {
	SqliteConn *sql.DB
	Mu         *sync.Mutex
}

func Init() *LiteDB {
	var db LiteDB
	db.Mu = &sync.Mutex{}
	db.Mu.Lock()
	fmt.Println("Got Init LiteDB lock")
	defer db.Mu.Unlock()
	defer fmt.Println("Rem Init LiteDB lock")
	conn, err := sql.Open("sqlite", "../tmp/gnuplex.sqlite3")
	db.SqliteConn = conn
	if err != nil {
		log.Fatal("Init LiteDB fatal error:", err)
	}
	_, err = db.SqliteConn.Exec("create table if not exists pos_cache (filepath string not null primary key, pos int);")
	if err != nil {
		log.Fatal("Init LiteDB error 1:", err)
	}
	_, err = db.SqliteConn.Exec("create table if not exists history (id integer not null unique, mediafile	text, primary key(id AUTOINCREMENT));")
	if err != nil {
		log.Fatal("Init LiteDB error 2:", err)
	}
	_, err = db.SqliteConn.Exec("create table if not exists medialist (filepath text not null,  primary key(filepath)) ;")
	if err != nil {
		log.Fatal("Init LiteDB error 3:", err)
	}
	_, err = db.SqliteConn.Exec("create table if not exists mediadirs (filepath text not null, primary key(filepath)) ;")
	if err != nil {
		log.Fatal("Init LiteDB error 4:", err)
	}
	_, err = db.SqliteConn.Exec("create table if not exists version_info (key string not null primary key, value string);")
	if err != nil {
		log.Fatal("Init LiteDB error 5:", err)
	}
	_, err = db.SqliteConn.Exec("insert or ignore into version_info values ('db_schema_version', ?);", consts.DBVersion)
	if err != nil {
		log.Fatal("Init LiteDB error 6:", err)
	}
	UpgradeDB(&db)
	return &db
}

func UpgradeDB(liteDB *LiteDB) {
	rows, err := liteDB.SqliteConn.Query("select value from version_info where key = 'db_schema_version';")
	if err != nil {
		log.Fatal("Upgrade db error:", err)
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
}
