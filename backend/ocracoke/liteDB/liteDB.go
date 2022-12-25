package liteDB

import (
	"database/sql"
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

func Init(prod bool) (*LiteDB, error) {
	var db LiteDB
	var conn *sql.DB
	var err error
	db.Mu = &sync.Mutex{}
	db.Mu.Lock()
	log.Println("Got Init LiteDB lock")
	defer db.Mu.Unlock()
	defer log.Println("Rem Init LiteDB lock")
	if prod {
		conn, err = sql.Open("sqlite", consts.ProdDBFilepath)
	} else {
		conn, err = sql.Open("sqlite", consts.DevDBFilepath)
	}
	db.SqliteConn = conn
	if err != nil {
		log.Println("Init LiteDB fatal error:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists pos_cache (filepath string not null primary key, pos int);")
	if err != nil {
		log.Println("Init LiteDB error 1:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists history (id integer not null unique, mediafile	text, primary key(id AUTOINCREMENT));")
	if err != nil {
		log.Println("Init LiteDB error 2:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists medialist (filepath text not null,  primary key(filepath)) ;")
	if err != nil {
		log.Println("Init LiteDB error 3:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists mediadirs (filepath text not null, primary key(filepath)) ;")
	if err != nil {
		log.Println("Init LiteDB error 4:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists file_exts (ext text not null, exclude int, primary key(ext)) ;")
	if err != nil {
		log.Println("Init LiteDB error 7:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("create table if not exists version_info (key string not null primary key, value string);")
	if err != nil {
		log.Println("Init LiteDB error 5:", err)
		return nil, err
	}
	_, err = db.SqliteConn.Exec("insert or ignore into version_info values ('db_schema_version', ?);", consts.DBVersion)
	if err != nil {
		log.Println("Init LiteDB error 6:", err)
		return nil, err
	}
	upgradeDB(&db)
	return &db, nil
}

func upgradeDB(db *LiteDB) error {
	rows, err := db.SqliteConn.Query("select value from version_info where key = 'db_schema_version';")
	if err != nil {
		log.Println("Upgrade db error:", err)
		return err
	}
	next := rows.Next()
	if !next {
		log.Fatal("Upgrade db error: no version schema", err)
	}
	var vers string
	err = rows.Scan(&vers)
	if err != nil {
		log.Println("Bad version schema 1", err)
		return err
	}
	versNum, err := strconv.Atoi(vers)
	if err != nil {
		log.Println("Bad version schema 2", err)
		return err
	}
	if versNum < consts.DBVersion {
		log.Println("Should I do something?")
	}
	rows.Close()
	return nil
}
