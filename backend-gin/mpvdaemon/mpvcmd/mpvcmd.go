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

func unixMsg(mpvConn *net.UnixConn, msg []byte) []byte {
	mpvConn.Write(append(msg, '\n'))
	readline := make([]byte, 2048)
	n, err := mpvConn.Read(readline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mpv cmd err", err)
	}
	return readline[:len(readline[:n])]
}

func mpvGetCmd(mpvConn *net.UnixConn, cmd []string) []byte {
	query := IMPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(mpvConn, jsonData)
}

func mpvSetCmd(mpvConn *net.UnixConn, cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(mpvConn, jsonData)
}

/*
 * MPV command public fxns
 */
func Play(mpvConn *net.UnixConn) []byte {
	return mpvSetCmd(mpvConn, []interface{}{"set_property", "pause", false})
}

func Pause(mpvConn *net.UnixConn) []byte {
	return mpvSetCmd(mpvConn, []interface{}{"set_property", "pause", true})
}

func IsPaused(mpvConn *net.UnixConn) []byte {
	return mpvGetCmd(mpvConn, []string{"get_property", "pause"})
}

func GetMedia(mpvConn *net.UnixConn) []byte {
	return mpvGetCmd(mpvConn, []string{"get_property", "path"})
}

func SetMedia(mpvConn *net.UnixConn, filepath string) []byte {
	return mpvSetCmd(mpvConn, []interface{}{"loadfile", filepath})
}

func GetVolume(mpvConn *net.UnixConn) []byte {
	return mpvGetCmd(mpvConn, []string{"get_property", "volume"})
}

func SetVolume(mpvConn *net.UnixConn, vol int) []byte {
	return mpvSetCmd(mpvConn, []interface{}{"set_property", "volume", vol})
}

func GetPos(mpvConn *net.UnixConn) []byte {
	return mpvGetCmd(mpvConn, []string{"get_property", "time-pos"})
}

func SetPos(mpvConn *net.UnixConn, pos int) []byte {
	return mpvSetCmd(mpvConn, []interface{}{"set_property", "time-pos", pos})
}
