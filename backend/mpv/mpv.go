package mpv

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

type MPV struct {
	Conn *net.UnixConn
	Mu   *sync.Mutex
}

func runDaemon(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocketPath string) {
	if !createMpvDaemon {
		return
	}
	defer wg.Done()
	for {
		var cmd *exec.Cmd
		if !verbose {
			cmd = exec.Command("mpv", "--idle=yes", fmt.Sprintf("--input-ipc-server=%s", mpvSocketPath), "--fs", "--save-position-on-quit")
		} else {
			cmd = exec.Command("mpv", "--idle=yes", "-v", fmt.Sprintf("--input-ipc-server=%s", mpvSocketPath), "--fs", "--save-position-on-quit")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Error: mpvdaemon.Run: ", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket string) (*MPV, error) {
	go runDaemon(wg, verbose, createMpvDaemon, mpvSocket)
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", mpvSocket)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var mpvConn *net.UnixConn
	for i := 10; mpvConn == nil || i >= 0; i-- {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(3 * time.Second)
		}
	}
	var mpv MPV
	mpv.Conn = mpvConn
	mpv.Mu = &sync.Mutex{}
	return &mpv, nil
}
