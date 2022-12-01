package mpvcmd

import "net"

func Play(mpvConn *net.UnixConn) {
	mpvConn.Write([]byte("{\"command\":[\"set_property\",\"pause\",false]}\n"))
}

func Pause(mpvConn *net.UnixConn) {
	mpvConn.Write([]byte("{\"command\":[\"set_property\",\"pause\",true]}\n"))
}
