package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"gnuplex-backend/consts"
	"gnuplex-backend/liteDB"
	"gnuplex-backend/mpvdaemon/mpvcmd"

	"github.com/gin-gonic/gin"
)

func Run(wg *sync.WaitGroup, db *liteDB.LiteDB) {
	defer wg.Done()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	mpvcmd.InitUnixConn()
	/*
	 * Serve static files
	 */
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	router.Static("/home", consts.StaticFilespath)
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
			AddHist(db, mediafile)
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
	router.GET("/api/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetMediadirs(db, false))
	})
	router.POST("/api/mediadirs", func(c *gin.Context) {
		mediadirsJson := []byte(c.Query("mediadirs"))
		var mediadirs []string
		err := json.Unmarshal(mediadirsJson, &mediadirs)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = SetMediadirs(db, mediadirs)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the mediadirs")
			}
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
		c.JSON(http.StatusOK, Last25(db))
	})
	router.GET("/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetMedialib(db, false))
	})
	router.POST("/api/medialist", func(c *gin.Context) {
		ScanLib(db)
		c.String(http.StatusOK, "OK")
	})
	/*
	 * Execution
	 */
	router.Run(":40000")
}
