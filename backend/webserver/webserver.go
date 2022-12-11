package webserver

import (
	"database/sql"
	"net/http"
	"strconv"
	"sync"

	"gnuplex-backend/mpvdaemon/mpvcmd"
	"gnuplex-backend/sqliteconn"

	"github.com/gin-gonic/gin"
)

func Run(wg *sync.WaitGroup, db *sql.DB) {
	defer wg.Done()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	mpvcmd.InitUnixConn()
	/*
	 * Serve static files
	 */
	router.Static("/gnuplex", "./public")
	/*
	 * API endpoints
	 */
	router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Play())
	})
	router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Pause())
	})
	router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.IsPaused())
	})
	router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetMedia())
	})
	router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			sqliteconn.AddHist(db, mediafile)
			c.Data(http.StatusOK, "application/json", mpvcmd.SetMedia(mediafile))
		}
	})
	router.GET("/api/vol", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetVolume())
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
			c.Data(http.StatusOK, "application/json", mpvcmd.SetVolume(vol))
		}
	})
	router.GET("/api/pos", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetPos())
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
			c.Data(http.StatusOK, "application/json", mpvcmd.SetPos(pos))
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
	router.Run(":40000")
}
