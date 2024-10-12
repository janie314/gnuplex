package main

import (
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
	noCreateMpvDaemon := flag.Bool("no_mpv_daemon", false, "Do not spawn an mpv daemon and mpv socket.")
	mpvSocket := flag.String("mpv_socket_path", "/tmp/mpvsocket", "Spawn an mpv daemon. Otherwise, use someone else's mpv socket.")
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
	oc, err := ocracoke.Init(&wg, *prod, *mpvSocket)
	if err != nil {
		log.Fatal(err)
	}
	go oc.Run(&wg)
	go mpvdaemon.Run(&wg, *verbose, !*noCreateMpvDaemon, *mpvSocket)
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
