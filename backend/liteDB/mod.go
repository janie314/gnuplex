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

func New(prod bool) (*LiteDB, error) {
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
	/*
	 * DB initialization SQL statements
	 */
	for i, statement := range consts.InitStatements() {
		_, err = db.SqliteConn.Exec(statement)
		if err != nil {
			log.Println("Init LiteDB error ", i, ":", err)
			return nil, err
		}
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
