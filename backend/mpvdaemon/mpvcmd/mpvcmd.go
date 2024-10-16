package mpvcmd

import (
	"bufio"
	"encoding/json"
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

var mu sync.Mutex
var mpvConn *net.UnixConn

func InitUnixConn(wg *sync.WaitGroup, mpvSocket string) {
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", mpvSocket)
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

/*
 * MPV command public fxns
 */
func Play() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", false})
}

func Pause() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", true})
}

func IsPaused() []byte {
	return mpvGetCmd([]string{"get_property", "pause"})
}

func GetMedia() []byte {
	return mpvGetCmd([]string{"get_property", "path"})
}

func SetMedia(filepath string) []byte {
	return mpvSetCmd([]interface{}{"loadfile", filepath})
}

func GetVolume() []byte {
	return mpvGetCmd([]string{"get_property", "volume"})
}

func SetVolume(vol int) []byte {
	return mpvSetCmd([]interface{}{"set_property", "volume", vol})
}

func GetPos() []byte {
	return mpvGetCmd([]string{"get_property", "time-pos"})
}

func SetPos(pos int) []byte {
	return mpvSetCmd([]interface{}{"set_property", "time-pos", pos})
}
