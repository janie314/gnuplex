package main

import (
	"context"
	"flag"
	"fmt"
	"gnuplex-backend/consts"
	"gnuplex-backend/mpvdaemon"
	"gnuplex-backend/ocracoke"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/quartz"
)

func main() {
	fmt.Println("GNUPlex Version " + consts.GNUPlexVersion)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched := quartz.NewStdScheduler()
	sched.Start(ctx)
	scanLibTrigger, err := quartz.NewCronTrigger("0 15 10 * * ?")
	if err != nil {
		log.Fatal("CronTrigger init failure", err)
	}
	scanLibJob := quartz.NewFunctionJob(func(_ context.Context) (int, error) {
		return 0, oc.ScanLib()
	})
	err = sched.ScheduleJob(ctx, scanLibJob, scanLibTrigger)
	if err != nil {
		log.Fatal("Scheduler init failure", err)
	}
	/*
	 * Main execution
	 */
	wg.Wait()
}
