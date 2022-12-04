package main

import (
	"gnuplex-backend/db"
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/webserver"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go webserver.Run(&wg)
	go db.Run(&wg)
	go mpvdaemon.Run(&wg)
	wg.Wait()
}
