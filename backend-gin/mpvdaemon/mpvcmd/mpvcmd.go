package mpvcmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func Play(mpvConn *net.UnixConn) string {
	mpvConn.Write([]byte(`{"command":["set_property","pause",false]}` + "\n"))
	readline := make([]byte, 1024)
	mpvConn.Read(readline)
	return strings.Split(string(readline), "\n")[0]
}

func Pause(mpvConn *net.UnixConn) string {
	mpvConn.Write([]byte(`{"command":["set_property","pause",true]}` + "\n"))
	readline := make([]byte, 1024)
	mpvConn.Read(readline)
	return strings.Split(string(readline), "\n")[0]
}

func LoadMedia(mpvConn *net.UnixConn, filepath string) string {
	data := map[string](interface{}){"command": []string{"loadmedia", filepath}}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	mpvConn.Write(append(jsonData, '\n'))
	readline := make([]byte, 1024)
	mpvConn.Read(readline)
	return strings.Split(string(readline), "\n")[0]
}
