package mpv

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// TODO: should the daemon be part of this struct
type MPV struct {
	mu   *sync.Mutex
	conn *net.UnixConn
}

/*
 * The MPV object.
 *
 * This objects needs a sync.WaitGroup; it is assumed that another
 * object is managing this object via a WaitGroups. This object
 * will control its video player daemon using this wg.
 */
func New(wg *sync.WaitGroup) (*MPV, error) {
	go RunDaemon(wg, false)
	var mpv_socket *net.UnixAddr
	mpv := new(MPV)
	mpv.mu = new(sync.Mutex)
	var conn *net.UnixConn
	mpv_socket, err := net.ResolveUnixAddr("unix", "/tmp/mpvsocket")
	if err != nil {
		return nil, err
	}
	for i := 10; conn == nil || i >= 0; i-- {
		conn, err = net.DialUnix("unix", nil, mpv_socket)
		if err != nil {
			log.Println("Warning: InitUnixConn:", err)
			time.Sleep(3 * time.Second)
		}
	}
	if err != nil {
		return nil, errors.New("couldn't get mpv Unix socket opened")
	} else {
		mpv.conn = conn
		return mpv, nil
	}
}
