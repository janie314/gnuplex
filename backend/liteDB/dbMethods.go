package liteDB

import (
	"log"
)

func (db *LiteDB) AddHist(mediafile string, ignorelock bool) error {
	db.Lock("AddHist", ignorelock)
	defer db.Unlock("AddHist", ignorelock)
	_, err := db.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddHist err", err)
	}
	return err
}
