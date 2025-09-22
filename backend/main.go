package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"gnuplex/consts"
	"gnuplex/gnuplex"
	server "gnuplex/gnuplex"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

// Set by Ldflag at compile time (see make.py's go_build task)
var SourceHash string
var Platform string
var GoVersion string

func main() {
	/*
	 * Cmd line flags
	 */
	prod := flag.Bool("prod", false, "Run in prod mode.")
	verbose := flag.Bool("verbose", false, "Verbose logging.")
	version := flag.Bool("version", false, "Print version info.")
	dbPath := flag.String("db_path", "gnuplex.sqlite3", "Path to sqlite DB.")
	port := flag.Int("port", 40000, "HTTP server's port.")
	upgrade := flag.Bool("upgrade", false, "Upgrade GNUPlex.")
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	staticFiles := flag.String("static_files", "static", "Path to static web files.")
	flag.Parse()
	// Some flags that subvert the main daemon process
	if *upgrade {
		if _, err := gnuplex.UpgradeGNUPlex(exe, true); err != nil {
			log.Fatalf("7a7233a9-262a-4bf6-8229-43855d3852d2 could not upgrade GNUPlex")
		}
		os.Exit(0)
	}
	if *version {
		printVersion()
	}
	fmt.Println("GNUPlex Version " + consts.Version)
	if *prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// Main daemon setup
	var wg sync.WaitGroup
	wg.Add(1)
	server, err := server.Init(&wg, (!*prod) || (*verbose), *dbPath, *staticFiles, *port, SourceHash, Platform, GoVersion, exe)
	if err != nil {
		log.Fatal(err)
	}
	go server.Run(&wg)
	// Scheduler process
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched, err := quartz.NewStdScheduler()
	if err != nil {
		log.Fatal("c98500e1-42f4-4c5d-ad2c-cedd4e4712b0 Failed to initialize cron scheduler", err)
	}
	sched.Start(ctx)
	scanLibTrigger, err := quartz.NewCronTrigger("0 15 1 * * *")
	if err != nil {
		log.Fatalln("9d9da752-4415-48ce-beec-0d8c703dd012 Failed to initialize cron scheduler", err)
	}
	scanLibJob := job.NewFunctionJob(func(_ context.Context) (int, error) {
		log.Println("running ScanLib job")
		return 0, server.ScanLib()
	})
	err = sched.ScheduleJob(quartz.NewJobDetail(scanLibJob, quartz.NewJobKey("scanlib")), scanLibTrigger)
	if err != nil {
		log.Fatalln("638eded7-2ad6-45b5-a13f-a99ad4642ff5 Failed to initialize cron scheduler", err)
	}
	updateTrigger, err := quartz.NewCronTrigger("0 15 2 * * *")
	if err != nil {
		log.Fatalln("2c007cd2-c363-44a5-81dd-6701a5487c9f Failed to initialize cron scheduler", err)
	}
	updateJob := job.NewFunctionJob(func(_ context.Context) (int, error) {
		log.Println("running update job")
		path, err := server.GetNowPlaying()
		if err != nil {
			log.Println("857034c0-a7db-4aa8-8cdf-00d0b6d811c2 Failed to retrive NowPlaying")
			return 0, err
		}
		if path != nil {
			log.Println("Not updating; something is playing")
		}
		upgraded, err := gnuplex.UpgradeGNUPlex(exe, false)
		if err != nil {
			return 0, err
		}
		if upgraded {
			log.Println("Updated GNUPlex. Quitting")
			server.Wg.Done()
		}
		return 0, nil
	})
	err = sched.ScheduleJob(quartz.NewJobDetail(updateJob, quartz.NewJobKey("update")), updateTrigger)
	if err != nil {
		log.Fatalln("6b379896-04c0-47e4-8655-a23f1126602a Failed to initialize cron scheduler", err)
	}
	// Main execution
	wg.Wait()
}

func printVersion() {
	var version consts.VersionInfo
	version.SourceHash = SourceHash
	version.Version = consts.Version
	version.Platform = Platform
	version.GoVersion = GoVersion
	res, err := json.Marshal(version)
	if err != nil {
		log.Fatalln("678f5d62-8c22-42bc-b25a-c5903b533312 failed to turn version info into JSON")
	}
	fmt.Println(string(res))
	os.Exit(0)
}
