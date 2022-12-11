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

func InitUnixConn() {
	var mpvUnixAddr *net.UnixAddr
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", "/tmp/mpvsocket")
	if err != nil {
		log.Fatal(err)
	}
	for mpvConn == nil {
		mpvConn, err = net.DialUnix("unix", nil, mpvUnixAddr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			time.Sleep(3 * time.Second)
		}
	}
}

func unixMsg(msg []byte) []byte {
	fmt.Println("waiting4lock", string(msg))
	mu.Lock()
	fmt.Println("GOTLOCK", string(msg))
	defer mu.Unlock()
	n, err := mpvConn.Write(append(msg, '\n'))
	fmt.Println("n1, err", n, err, string(msg))
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
	fmt.Println("J", string(jsonData))
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
