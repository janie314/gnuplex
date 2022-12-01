package mpvcmd

import (
	"net"
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
