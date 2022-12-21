package mpvdaemon

import (
	"log"
	"os/exec"
	"sync"
	"time"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		cmd := exec.Command("mpv", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
		err := cmd.Run()
		if err != nil {
			log.Println("Error: mpvdaemon.Run: ", err)
		}
		time.Sleep(3 * time.Second)
	}
}
