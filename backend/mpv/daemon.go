package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func Run(wg *sync.WaitGroup, verbose bool) {
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
		time.Sleep(3 * time.Second)
	}
}

/*
 * Types
 */
type IMPVQuery struct {
	Command []interface{} `json:"command"`
}

type IMPVQueryString struct {
	Command []string `json:"command"`
}

type IMPVResponseBool struct {
	Data bool `json:"data"`
}

type IMPVResponseString struct {
	Data string `json:"data"`
}

type IMPVResponseInt struct {
	Data int `json:"data"`
}

/*
 * Aux fxns
 */

var mu sync.Mutex
var mpvConn *net.UnixConn

func InitUnixConn(wg *sync.WaitGroup) {
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", "/tmp/mpvsocket")
	if err != nil {
		log.Fatal(err)
	}
	for i := 10; mpvConn == nil || i >= 0; i-- {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(3 * time.Second)
		}
	}
}

func unixMsg(msg []byte) []byte {
	mu.Lock()
	defer mu.Unlock()
	_, err := mpvConn.Write(append(msg, '\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(mpvConn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "request_id") {
			return []byte(line)
		}
	}
	return []byte{}
}

func mpvGetCmd(cmd []string) []byte {
	query := IMPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(jsonData)
}

func mpvSetCmd(cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(jsonData)
}
