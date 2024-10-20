package mpv

import (
	"bufio"
	"encoding/json"
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

type MPVResult[T string | int | float64] struct {
	Data      T      `json:"data"`
	RequestId int    `json:"request_id"`
	Error     string `json:"error"`
}

func processMpvResult[T string | int | float64](resBytes []byte) (T, error) {
	var res MPVResult[T]
	var defaultVal T
	err := json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("mpv result error", err)
		return defaultVal, err
	} else if res.Error != "success" {
		log.Println("mpv result error", err)
		return defaultVal, errors.New(res.Error)
	}
	return res.Data, nil

}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket string) (*net.UnixConn, error) {
	go RunDaemon(wg, verbose, createMpvDaemon, mpvSocket)
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", mpvSocket)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var mpvConn *net.UnixConn
	for i := 10; mpvConn == nil || i >= 0; i-- {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(3 * time.Second)
		}
	}
	return mpvConn, nil
}

func unixMsg(mpvConn *net.UnixConn, mu *sync.Mutex, msg []byte) []byte {
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

func mpvGetCmd(mpvConn *net.UnixConn, mu *sync.Mutex, cmd []string) []byte {
	query := IMPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(mpvConn, mu, jsonData)
}

func mpvSetCmd(mpvConn *net.UnixConn, mu *sync.Mutex, cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(mpvConn, mu, jsonData)
}

/*
 * MPV command public fxns
 */
func Play(mpvConn *net.UnixConn, mu *sync.Mutex) []byte {
	return mpvSetCmd(mpvConn, mu, []interface{}{"set_property", "pause", false})
}

func Pause(mpvConn *net.UnixConn, mu *sync.Mutex) []byte {
	return mpvSetCmd(mpvConn, mu, []interface{}{"set_property", "pause", true})
}

func IsPaused(mpvConn *net.UnixConn, mu *sync.Mutex) []byte {
	return mpvGetCmd(mpvConn, mu, []string{"get_property", "pause"})
}

func GetMedia(mpvConn *net.UnixConn, mu *sync.Mutex) []byte {
	return mpvGetCmd(mpvConn, mu, []string{"get_property", "path"})
}

func SetMedia(mpvConn *net.UnixConn, mu *sync.Mutex, filepath string) []byte {
	return mpvSetCmd(mpvConn, mu, []interface{}{"loadfile", filepath})
}

func GetVol(mpvConn *net.UnixConn, mu *sync.Mutex) (int, error) {
	resBytes := mpvGetCmd(mpvConn, mu, []string{"get_property", "volume"})
	n, err := processMpvResult[float64](resBytes)
	if err != nil {
		return 0, err
	}
	return int(n), err
}

func SetVolume(mpvConn *net.UnixConn, mu *sync.Mutex, vol int) []byte {
	return mpvSetCmd(mpvConn, mu, []interface{}{"set_property", "volume", vol})
}

func GetPos(mpvConn *net.UnixConn, mu *sync.Mutex) (int, error) {
	resBytes := mpvGetCmd(mpvConn, mu, []string{"get_property", "time-pos"})
	n, err := processMpvResult[float64](resBytes)
	if err != nil {
		return 0, err
	}
	return int(n), err
}

func GetTimeRemaining(mpvConn *net.UnixConn, mu *sync.Mutex) []byte {
	return mpvGetCmd(mpvConn, mu, []string{"get_property", "time-remaining"})
}

func SetPos(mpvConn *net.UnixConn, mu *sync.Mutex, pos int) []byte {
	return mpvSetCmd(mpvConn, mu, []interface{}{"set_property", "time-pos", pos})
}
