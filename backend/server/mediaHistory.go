package server

import (
	"fmt"
	"log"
	"os"
)

func (oc *Server) Last25() []string {
	oc.DB.Mu.Lock()
	log.Println("Got Last25 lock")
	defer oc.DB.Mu.Unlock()
	defer fmt.Println("Rem Last25 lock")
	rows, err := oc.DB.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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

func (server *Server) AddHist(mediafile string) error {
	server.DB.Mu.Lock()
	log.Println("Got AddHist lock")
	defer server.DB.Mu.Unlock()
	defer log.Println("Rem AddHist lock")
	_, err := server.DB.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddHist err", err)
	}
	return err
}
