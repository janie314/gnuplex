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
	RestartCount int
	Verbose      bool
}

func (mpv *MPV) restartProcess() error {
	if mpv.RestartCount >= 10 {
		return errors.New("4b367530-abab-4964-a093-72f78cf523f2 reached limit of MPV restarts (10)")
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
		log.Println("Error 1b2b2f70-d75c-4ee4-b977-95f7272f72e9: mpvdaemon.Run: ", err)
		return err
	}
	// Initialize Unix socket
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", mpvSocketPath)
	if err != nil {
		return err
	}
	var mpvConn *net.UnixConn
	for i := 100; (err != nil || i == 100) && i >= 0; i-- {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			log.Println("Warning 6b94b8d9-9992-4276-8f6f-4d1b339cb59a: could not init unix socket to mpv, trying again in 2 seconds.", err)
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		return fmt.Errorf("84c2ca00-1df2-4af5-8bd4-9248d941ad7d: could not establish a unix socket connection to %s", mpvSocketPath)
	} else {
		log.Println("Established unix socket connection to mpv.")
		mpv.Conn = mpvConn
	}
	return nil
}

func (mpv *MPV) runProcessSupervisor(wg *sync.WaitGroup) {
	for {
		if _, err := mpv.Ping(); err != nil {
			log.Println("Error b304c2bb-a940-4a3e-979d-2939e986ed87: failed to ping MPV. restarting MPV.")
			if err := mpv.restartProcess(); err != nil {
				log.Println("Error c7b92229-b6b8-4800-91d6-8b8afb733d9c: failed to restart MPV. stopping gnuplex.")
				wg.Done()
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func Init(wg *sync.WaitGroup, verbose bool) (*MPV, error) {
	var mpv MPV
	mpv.Mu = &sync.Mutex{}
	mpv.RestartCount = 0
	mpv.Verbose = verbose
	if err := mpv.restartProcess(); err != nil {
		return nil, err
	}
	go mpv.runProcessSupervisor(wg)
	return &mpv, nil
}
