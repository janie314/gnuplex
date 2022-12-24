package mpvdaemon

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func Run(wg *sync.WaitGroup, debug bool) {
	defer wg.Done()
	for {
		var cmd *exec.Cmd
		if !debug {
			cmd = exec.Command("mpv", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
		} else {
			cmd = exec.Command("mpv", "--idle=yes", "-v", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
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
