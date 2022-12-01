package webserver

import (
	"net"
	"net/http"
	"sync"

	"gnuplex-backend/mpvdaemon/mpvcmd"

	"github.com/gin-gonic/gin"
)

func Run(wg *sync.WaitGroup, mpvConn *net.UnixConn) {
	defer wg.Done()
	router := gin.Default()
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
