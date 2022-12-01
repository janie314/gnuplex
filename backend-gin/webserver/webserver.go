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
	router.POST("/play", func(c *gin.Context) {
		mpvcmd.Play(mpvConn)
		c.JSON(http.StatusOK, "OK")
	})
	router.POST("/pause", func(c *gin.Context) {
		mpvcmd.Pause(mpvConn)
		c.JSON(http.StatusOK, "OK")
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
