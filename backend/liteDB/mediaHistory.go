package liteDB

import (
	"fmt"
	"log"
	"os"
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

func (db *LiteDB) Last25(ignorelock bool) []string {
	db.Lock("Last25", ignorelock)
	defer db.Unlock("Last25", ignorelock)
	rows, err := db.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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
			fmt.Println("Error: Last25:", err)
		} else if i < len(res) {
			res[i] = str
			i++
		} else {
			res = append(res, str)
		}
	}
	return res[:i]
}
