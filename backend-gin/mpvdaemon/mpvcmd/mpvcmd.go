package mpvcmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type IMPVResponseBool struct {
	Data bool `json:"data"`
}

func mpvAuxFxn(mpvConn *net.UnixConn, cmd string) []byte {
	mpvConn.Write([]byte(cmd + "\n"))
	readline := make([]byte, 1024)
	n, err := mpvConn.Read(readline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mpv cmd err", err)
	}
	return readline[:len(readline[:n])]
}

func Play(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["set_property","pause",false]}`)
}

func Pause(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["set_property","pause",true]}`)
}

func IsPaused(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["get_property","pause"]}`)
}

func GetMedia(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["get_property","path"]}`)
}

func GetVolume(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["get_property","volume"]}`)
}

func GetPos(mpvConn *net.UnixConn) []byte {
	return mpvAuxFxn(mpvConn, `{"command":["get_property","time-pos"]}`)
}

func LoadMedia(mpvConn *net.UnixConn, filepath string) []byte {
	data := map[string](interface{}){"command": []string{"loadmedia", filepath}}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loadmedia problem", err)
		return []byte{}
	}
	mpvConn.Write(append(jsonData, '\n'))
	readline := make([]byte, 1024)
	mpvConn.Read(readline)
	return readline
}
