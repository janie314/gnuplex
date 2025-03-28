package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"gnuplex/consts"
	server "gnuplex/gnuplex"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

// Set by Ldflag at compile time (see Rakefile's go_build task)
var SourceHash string

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
		upgradeGNUPlex(exe)
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
	server, err := server.Init(&wg, (!*prod) || (*verbose), *dbPath, *staticFiles, *port, SourceHash)
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
	scanLibTrigger, err := quartz.NewCronTrigger("0 15 10 * * ?")
	if err != nil {
		log.Fatalln("9d9da752-4415-48ce-beec-0d8c703dd012 Failed to initialize cron scheduler", err)
	}
	scanLibJob := job.NewFunctionJob(func(_ context.Context) (int, error) {
		return 0, server.ScanLib()
	})
	err = sched.ScheduleJob(quartz.NewJobDetail(scanLibJob, quartz.NewJobKey("scanlib")), scanLibTrigger)
	if err != nil {
		log.Fatalln("638eded7-2ad6-45b5-a13f-a99ad4642ff5 Failed to initialize cron scheduler", err)
	}
	// Main execution
	wg.Wait()
}

func upgradeGNUPlex(exe string) {
	cmd := exec.Command("git", "-C", filepath.Join(filepath.Dir(exe), "../.."), "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("f0ac0db2-c77e-4bb4-94b4-7b98931d6379 Failed to upgrade GNUPlex", err)
	} else {
		fmt.Println("Successfully upgraded! Now run `systemctl --user restart gnuplex`.")
		os.Exit(0)
	}
}

func printVersion() {
	var version consts.VersionInfo
	version.SourceHash = SourceHash
	version.Version = consts.Version
	res, err := json.Marshal(version)
	if err != nil {
		log.Fatalln("678f5d62-8c22-42bc-b25a-c5903b533312 failed to turn version info into JSON")
	}
	fmt.Println(string(res))
	os.Exit(0)
}
