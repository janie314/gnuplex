package webserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"gnuplex-backend/mpvdaemon/mpvcmd"

	"github.com/gin-gonic/gin"
)

func initUnixConn() *net.UnixConn {
	var mpvUnixAddr *net.UnixAddr
	var mpvConn *net.UnixConn
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
	return mpvConn
}

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	router := gin.Default()
	mpvConn := initUnixConn()
	/*
	 * API endpoints
	 */
	router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(mpvcmd.Play(mpvConn)))
	})
	router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(mpvcmd.Pause(mpvConn)))
	})
	/*
	 * Serve static files
	 */
	router.Static("/", "./public")
	/*
	 * Execution
	 */
	router.Run(":50000")
}
