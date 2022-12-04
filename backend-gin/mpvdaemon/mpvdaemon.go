package mpvdaemon

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("mpv", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--fs", "--save-position-on-quit")
	err := cmd.Run()
	fmt.Fprintln(os.Stderr, err)
}
