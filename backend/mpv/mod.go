package mpv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

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

type MPV struct {
	mu   sync.Mutex
	conn *net.UnixConn
}

func New(wg *sync.WaitGroup) (*MPV, error) {
	go RunDaemon(wg, false)
	var mpv MPV
	var mpv_socket *net.UnixAddr
	mpv_socket, err := net.ResolveUnixAddr("unix", "/tmp/mpvsocket")
	if err != nil {
		return nil, err
	}
	for i := 10; mpv.conn == nil || i >= 0; i-- {
		mpv.conn, err = net.DialUnix("unix", nil, mpv_socket)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(3 * time.Second)
		}
	}
	if err != nil {
		return nil, errors.New("couldn't get mpv Unix socket opened")
	} else {
		return &mpv, nil
	}
}

func (mpv *MPV) UnixMsg(msg []byte) []byte {
	mpv.mu.Lock()
	defer mpv.mu.Unlock()
	_, err := mpv.conn.Write(append(msg, '\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(mpv.conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "request_id") {
			log.Println("debug", line)
			return []byte(line)
		}
	}
	return []byte{}
}
