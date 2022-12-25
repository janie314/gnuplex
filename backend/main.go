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
	prod := flag.Bool("prod", false, "Run in prod mode.")
	verbose := flag.Bool("verbose", false, "Verbose logging.")
	flag.Parse()
	if *prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	/*
	 * Main daemon setup
	 */
	var wg sync.WaitGroup
	wg.Add(1)
	oc, err := ocracoke.Init(&wg, *prod)
	if err != nil {
		log.Fatal(err)
	}
	go oc.Run(&wg)
	go mpvdaemon.Run(&wg, *verbose)
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
