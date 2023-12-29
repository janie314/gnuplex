package server

import (
	"fmt"
	"log"
	"os"
)

func (srv *Server) Last25() []string {
	srv.DB.Mu.Lock()
	log.Println("Got Last25 lock")
	defer srv.DB.Mu.Unlock()
	defer fmt.Println("Rem Last25 lock")
	rows, err := srv.DB.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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
