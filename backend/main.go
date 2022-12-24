package main

import (
	"flag"
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/ocracoke"
	"log"
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
	wg.Add(1)
	oc, err := ocracoke.Init(&wg, *debug)
	if err != nil {
		log.Fatal(err)
	}
	go oc.Run(&wg)
	go mpvdaemon.Run(&wg, *debug)
	/*
	 * Scheduler process
	 */
	sched := quartz.NewStdScheduler()
	sched.Start()
	scanLibTrigger, err := quartz.NewCronTrigger("0 15 10 * * ?")
	if err != nil {
		log.Fatal("CronTrigger init failure", err)
	}
	scanLibJob := quartz.NewFunctionJob(func() (int, error) {
		return 0, oc.ScanLib()
	})
	err = sched.ScheduleJob(scanLibJob, scanLibTrigger)
	if err != nil {
		log.Fatal("Scheduler init failure", err)
	}
	/*
	 * Main execution
	 */
	wg.Wait()
}
