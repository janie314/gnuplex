package main

import (
	"flag"
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/sqliteconn"
	"gnuplex-backend/webserver"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	debug := flag.Bool("debug", false, "Enable non-production mode + more verbose logging.")
	flag.Parse()
	if !(*debug) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	/*
	 * Execution
	 */
	var wg sync.WaitGroup
	wg.Add(2)
	db := sqliteconn.Init()
	go webserver.Run(&wg, db)
	go mpvdaemon.Run(&wg)
	wg.Wait()
}
