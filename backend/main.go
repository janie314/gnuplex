package main

import (
	"context"
	"flag"
	"fmt"
	"gnuplex-backend/consts"
	"gnuplex-backend/server"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/quartz"
)

func main() {
	fmt.Println("GNUPlex Server Version " + consts.GNUPlexVersion)
	/*
	 * Cmd line flags
	 */
	prod := flag.Bool("prod", false, "Run in prod mode.")
	port := flag.Int("port", 40000, "Port to listen on.")
	db := flag.String("db", "gnuplex.sqlite3", "Filepath to SQLite database.")
	api_url_base := flag.String("api_url_base", "/api", "Base URL for server HTTP requests")
	static_url_base := flag.String("static_url_base", "/home", "Base URL for static files")
	static_dir := flag.String("static_dir", "public", "Static file directory")
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
	srv, err := server.New(&wg, *prod, *port, *static_url_base, *api_url_base, *db, *static_dir)
	if err != nil {
		log.Fatal(err)
	}
	go srv.Run(&wg)
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
		return 0, srv.ScanLib(false)
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
