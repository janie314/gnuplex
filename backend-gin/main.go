package main

import (
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/sqliteconn"
	"gnuplex-backend/webserver"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	db := sqliteconn.Init()
	go webserver.Run(&wg, db)
	go mpvdaemon.Run(&wg)
	wg.Wait()
}
