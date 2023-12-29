package mpv

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func RunDaemon(wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()
	for {
		var cmd *exec.Cmd
		if !verbose {
			cmd = exec.Command("mpv", "--cursor-autohide=always", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
		} else {
			cmd = exec.Command("mpv", "--cursor-autohide=always", "--idle=yes", "-v", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Error: mpvdaemon.Run: ", err)
		}
	}
}
