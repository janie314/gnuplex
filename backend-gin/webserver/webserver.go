package webserver

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gnuplex-backend/mpvdaemon/mpvcmd"
	"gnuplex-backend/sqliteconn"

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

func Run(wg *sync.WaitGroup, db *sql.DB) {
	defer wg.Done()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	mpvConn := initUnixConn()
	/*
	 * Serve static files
	 */
	router.Static("/gnuplex", "./public")
	/*
	 * API endpoints
	 */
	router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Play(mpvConn))
	})
	router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Pause(mpvConn))
	})
	router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.IsPaused(mpvConn))
	})
	router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetMedia(mpvConn))
	})
	router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			c.Data(http.StatusOK, "application/json", mpvcmd.SetMedia(mpvConn, mediafile))
		}
	})
	router.GET("/api/vol", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetVolume(mpvConn))
	})
	router.POST("/api/vol", func(c *gin.Context) {
		param := c.Query("vol")
		if param == "" {
			c.String(http.StatusBadRequest, "empty vol string")
		} else {
			vol, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad vol string")
			}
			c.Data(http.StatusOK, "application/json", mpvcmd.SetVolume(mpvConn, vol))
		}
	})
	router.GET("/api/pos", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetPos(mpvConn))
	})
	router.POST("/api/pos", func(c *gin.Context) {
		param := c.Query("pos")
		if param == "" {
			c.String(http.StatusBadRequest, "empty pos string")
		} else {
			pos, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad pos string")
			}
			c.Data(http.StatusOK, "application/json", mpvcmd.SetPos(mpvConn, pos))
		}
	})
	router.GET("/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, sqliteconn.Last25(db))
	})
	router.GET("/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, sqliteconn.GetMedialib(db))
	})
	/*
	 * Execution
	 */
	router.Run(":50000")
}
