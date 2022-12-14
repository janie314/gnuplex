package main

import (
	"flag"
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/sqliteconn"
	"gnuplex-backend/webserver"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/quartz"
)

func main() {
	/*
	 * Cmd line flags
	 */
	debug := flag.Bool("debug", false, "Enable non-production mode + more verbose logging.")
	flag.Parse()
	if !(*debug) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	/*
	 * Main daemon setup
	 */
	var wg sync.WaitGroup
	wg.Add(2)
	db := sqliteconn.Init()
	go webserver.Run(&wg, db)
	go mpvdaemon.Run(&wg)
	/*
	 * Scheduler process
	 */
	sched := quartz.NewStdScheduler()
	sched.Start()
	scanLibTrigger, _ := quartz.NewCronTrigger("13 10 * * *")
	scanLibJob := quartz.NewFunctionJob(func() (int, error) { return 0, sqliteconn.ScanLib((db)) })
	sched.ScheduleJob(scanLibJob, scanLibTrigger)
	/*
	 * Main execution
	 */
	wg.Wait()
}
