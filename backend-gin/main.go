package main

import (
	"gnuplex-backend/db"
	"gnuplex-backend/webserver"
	"log"
	"net"
	"sync"
)

func main() {
	mpvUnixAddr, err := net.ResolveUnixAddr("unix", "/tmp/mpvsocket")
	if err != nil {
		log.Fatal(err)
	}
	mpvConn, err := net.DialUnix("unix", nil, mpvUnixAddr)
	if err != nil {
		log.Fatal(err)
	}
	/*
	 * Main execution
	 */
	var wg sync.WaitGroup
	wg.Add(2)
	go webserver.Run(&wg, mpvConn)
	go db.Run(&wg)
	wg.Wait()
}
