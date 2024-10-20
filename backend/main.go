package main

import (
	"context"
	"flag"
	"fmt"
	"gnuplex-backend/consts"
	server "gnuplex-backend/gnuplex"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/job"
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
	dbPath := flag.String("db_path", "gnuplex.sqlite3", "Path to sqlite DB.")
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	staticFiles := flag.String("static_files", filepath.Join(filepath.Dir(exe), "static"), "Path to static web files.")
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
	server, err := server.Init(&wg, (!*prod) || (*verbose), !*noCreateMpvDaemon, *mpvSocket, *dbPath, *staticFiles)
	if err != nil {
		log.Fatal(err)
	}
	go server.Run(&wg)
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
	scanLibJob := job.NewFunctionJob(func(_ context.Context) (int, error) {
		return 0, server.ScanLib()
	})
	err = sched.ScheduleJob(quartz.NewJobDetail(scanLibJob, quartz.NewJobKey("scanlib")), scanLibTrigger)
	if err != nil {
		log.Fatal("Scheduler init failure", err)
	}
	/*
	 * Main execution
	 */
	wg.Wait()
}
