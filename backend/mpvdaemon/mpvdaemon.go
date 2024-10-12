package mpvdaemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func Run(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocketPath string) {
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
