package gnuplex

import (
	"gnuplex/consts"
	"gnuplex/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaActionBody struct {
	Id models.MediaItemId `json:"id"`
}

type CastBody struct {
	Url  string `json:"url"`
	Temp bool   `json:"temp"`
}

type MediaDirsBody []string

type FileExtsBody []string

// Initialize the web server's HTTP Endpoints
func (gnuplex *GNUPlex) InitWebEndpoints(prod bool, staticFiles string) {
	gnuplex.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	gnuplex.Router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	gnuplex.Router.Static("/home", staticFiles)
	gnuplex.Router.GET("/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	gnuplex.Router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", gnuplex.MPV.Play())
	})
	gnuplex.Router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", gnuplex.MPV.Pause())
	})
	gnuplex.Router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", gnuplex.MPV.IsPaused())
	})
	gnuplex.Router.GET("/api/nowplaying", func(c *gin.Context) {
		media, err := gnuplex.MPV.GetNowPlaying()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, media)
		}
	})
	gnuplex.Router.POST("/api/nowplaying", func(c *gin.Context) {
		body := MediaActionBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.ReplaceQueueAndPlay(body.Id); err != nil {
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.POST("/api/cast", func(c *gin.Context) {
		body := CastBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.Cast(body.Url, body.Temp); err != nil {
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})

	gnuplex.Router.GET("/api/vol", func(c *gin.Context) {
		vol, err := gnuplex.MPV.GetVol()
		if err != nil {
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})
	gnuplex.Router.POST("/api/vol", func(c *gin.Context) {
		param := c.Query("vol")
		if param == "" {
			c.String(http.StatusBadRequest, "empty vol string")
		} else {
			vol, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad vol string")
			} else {
				c.Data(http.StatusOK, "application/json", gnuplex.MPV.SetVolume(vol))
			}
		}
	})
	gnuplex.Router.GET("/api/mediadirs", func(c *gin.Context) {
		res, err := gnuplex.DB.GetMediaDirs()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/mediadirs", func(c *gin.Context) {
		body := MediaDirsBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.DB.SetMediadirs(body); err != nil {
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.GET("/api/file_exts", func(c *gin.Context) {
		res, err := gnuplex.DB.GetFileExts()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/file_exts", func(c *gin.Context) {
		body := FileExtsBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.DB.SetFileExts(body); err != nil {
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})

	gnuplex.Router.GET("/api/pos", func(c *gin.Context) {
		vol, err := gnuplex.MPV.GetPos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})

	gnuplex.Router.GET("/api/timeremaining", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", gnuplex.MPV.GetTimeRemaining())
	})
	gnuplex.Router.POST("/api/pos", func(c *gin.Context) {
		param := c.Query("pos")
		if param == "" {
			c.String(http.StatusBadRequest, "empty pos string")
		} else {
			pos, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad pos string")
			} else {
				c.Data(http.StatusOK, "application/json", gnuplex.MPV.SetPos(pos))
			}
		}
	})
	gnuplex.Router.GET("/api/last25", func(c *gin.Context) {
		res, err := gnuplex.DB.GetLast25Played()
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}

	})
	gnuplex.Router.GET("/api/mediaitems", func(c *gin.Context) {
		res, err := gnuplex.DB.GetMediaItems()
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/scanlib", func(c *gin.Context) {
		if err := gnuplex.ScanLib(); err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})

}
