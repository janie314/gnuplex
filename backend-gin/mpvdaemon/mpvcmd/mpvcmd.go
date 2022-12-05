package mpvcmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

/*
 * Types
 */
type IMPVQueryString struct {
	Command []interface{} `json:"command"`
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
func mpvGet(mpvConn *net.UnixConn, cmd string) []byte {
	mpvConn.Write([]byte(cmd + "\n"))
	readline := make([]byte, 1024)
	n, err := mpvConn.Read(readline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mpv cmd err 1", err)
	}
	return readline[:len(readline[:n])]
}

func mpvSetString(mpvConn *net.UnixConn, cmd, val string) []byte {
	query := IMPVQueryString{Command: []interface{}{cmd, val}}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	mpvConn.Write(append(jsonData, '\n'))
	readline := make([]byte, 1024)
	n, err := mpvConn.Read(readline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mpv cmd err 2", err)
	}
	return readline[:len(readline[:n])]
}

/*
 * MPV command public fxns
 */
func Play(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["set_property","pause",false]}`)
}

func Pause(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["set_property","pause",true]}`)
}

func IsPaused(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["get_property","pause"]}`)
}

func GetMedia(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["get_property","path"]}`)
}

func SetMedia(mpvConn *net.UnixConn, filepath string) []byte {
	return mpvSetString(mpvConn, "loadmedia", filepath)
}

func GetVolume(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["get_property","volume"]}`)
}

func GetPos(mpvConn *net.UnixConn) []byte {
	return mpvGet(mpvConn, `{"command":["get_property","time-pos"]}`)
}
