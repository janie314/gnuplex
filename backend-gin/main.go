package main

import (
	"gnuplex-backend/webserver"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go webserver.Run(&wg)
	wg.Wait()
}
