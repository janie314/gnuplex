package mpvdaemon

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		cmd := exec.Command("mpv", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
		err := cmd.Run()
		fmt.Fprintln(os.Stderr, err)
		time.Sleep(3 * time.Second)
	}
}
