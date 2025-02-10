package mpv

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type MPV struct {
	Conn         *net.UnixConn
	Mu           *sync.Mutex
	Process      *os.Process
	Queue        []string
	RestartCount int
	Verbose      bool
}

func (mpv *MPV) restartProcess() error {
	if mpv.RestartCount >= 10 {
		return errors.New("reached restart limit, 10")
	}
	mpv.RestartCount += 1
	if mpv.Process != nil {
		if err := syscall.Kill(mpv.Process.Pid, syscall.Signal(0)); err == nil {
			if err := mpv.Process.Kill(); err != nil {
				return err
			}
		}
	}
	mpvSocketPath := path.Join("/tmp", fmt.Sprintf("mpvsocket-%s", uuid.New().String()))
	var cmd *exec.Cmd
	if !mpv.Verbose {
		cmd = exec.Command("mpv", "--idle=yes", fmt.Sprintf("--input-ipc-server=%s", mpvSocketPath), "--fs", "--save-position-on-quit")
	} else {
		cmd = exec.Command("mpv", "--idle=yes", "-v", fmt.Sprintf("--input-ipc-server=%s", mpvSocketPath), "--fs", "--save-position-on-quit")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("starting mpv process")
	if err := cmd.Start(); err != nil {
		log.Println("Error: mpvdaemon.Run: ", err)
		return err
	}
	// Initialize Unix socket
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", mpvSocketPath)
	if err != nil {
		return err
	}
	var mpvConn *net.UnixConn
	for i := 10; (err != nil || i == 10) && i >= 0; i-- {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		return fmt.Errorf("could not establish a unix socket connection to %s", mpvSocketPath)
	} else {
		log.Println("stablished unix socket connection")
		mpv.Conn = mpvConn
	}
	return nil
}

func (mpv *MPV) runProcessSupervisor(wg *sync.WaitGroup) {
	for {
		if _, err := mpv.Ping(); err != nil {
			if err := mpv.restartProcess(); err != nil {
				wg.Done()
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func Init(wg *sync.WaitGroup, verbose bool) (*MPV, error) {
	var mpv MPV
	mpv.Mu = &sync.Mutex{}
	mpv.Queue = make([]string, 0)
	mpv.RestartCount = 0
	mpv.Verbose = verbose
	if err := mpv.restartProcess(); err != nil {
		return nil, err
	}
	go mpv.runProcessSupervisor(wg)
	return &mpv, nil
}
