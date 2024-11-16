package main

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

var SourceHash string

func main() {
	/*
	 * Cmd line flags
	 */
	prod := flag.Bool("prod", false, "Run in prod mode.")
	verbose := flag.Bool("verbose", false, "Verbose logging.")
	version := flag.Bool("version", false, "Print version.")
	noCreateMpvDaemon := flag.Bool("no_mpv_daemon", false, "Do not spawn an mpv daemon and mpv socket.")
	mpvSocket := flag.String("mpv_socket_path", fmt.Sprintf("/tmp/mpvsocket-%s", uuid.New().String()), "Spawn an mpv daemon. Otherwise, use someone else's mpv socket.")
	dbPath := flag.String("db_path", "gnuplex.sqlite3", "Path to sqlite DB.")
	upgrade := flag.Bool("upgrade", false, "Upgrade GNUPlex.")
	source_hash := flag.Bool("source_hash", false, "Git commit this build comes from.")
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
	if *source_hash {
		sourceHash()
	}
	if *version {
		printVersion()
	}
	fmt.Println("GNUPlex Version " + consts.GNUPlexVersion)
	if *prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// Main daemon setup
	var wg sync.WaitGroup
	wg.Add(1)
	server, err := server.Init(&wg, (!*prod) || (*verbose), !*noCreateMpvDaemon, *mpvSocket, *dbPath, *staticFiles)
	if err != nil {
		log.Fatal(err)
	}
	go server.Run(&wg)
	// Scheduler process
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
	// Main execution
	wg.Wait()
}

func upgradeGNUPlex(exe string) {
	cmd := exec.Command("git", "-C", filepath.Join(filepath.Dir(exe), "../.."), "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("fail", err)
	} else {
		os.Exit(0)
	}
}

func sourceHash() {
	fmt.Println(SourceHash)
	os.Exit(0)
}

func printVersion() {
	fmt.Println(consts.GNUPlexVersion)
	os.Exit(0)
}
