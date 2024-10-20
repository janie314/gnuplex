package gnuplex

import (
	"encoding/json"
	"gnuplex-backend/consts"
	"gnuplex-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaActionBody struct {
	Id models.MediaItemId `json:"id"`
}

// Initialize the web server's HTTP Endpoints
func (gnuplex *GNUPlex) InitWebEndpoints(prod bool) {
	gnuplex.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	gnuplex.Router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	if prod {
		gnuplex.Router.Static("/home", consts.ProdStaticFilespath)
	} else {
		gnuplex.Router.Static("/home", consts.DevStaticFilespath)
	}
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
	gnuplex.Router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", gnuplex.MPV.GetMedia())
	})
	gnuplex.Router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			gnuplex.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", gnuplex.MPV.SetMedia(mediafile))
		}
	})
	gnuplex.Router.POST("/api/media2", func(c *gin.Context) {
		body := MediaActionBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.ReplaceQueueAndPlay(body.Id); err != nil {
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
			}
			c.Data(http.StatusOK, "application/json", gnuplex.MPV.SetVolume(vol))
		}
	})
	gnuplex.Router.GET("/api/mediadirs", func(c *gin.Context) {
		res, err := gnuplex.NewDB.GetMediaDirs()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/mediadirs", func(c *gin.Context) {
		mediadirsJson := []byte(c.Query("mediadirs"))
		var mediadirs []string
		err := json.Unmarshal(mediadirsJson, &mediadirs)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = gnuplex.SetMediadirs(mediadirs)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the mediadirs")
			}
		}
	})
	gnuplex.Router.GET("/api/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gnuplex.GetFileExts(false))
	})
	gnuplex.Router.POST("/api/file_exts", func(c *gin.Context) {
		fileExtsJson := []byte(c.Query("file_exts"))
		var fileExts []string
		err := json.Unmarshal(fileExtsJson, &fileExts)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = gnuplex.SetFileExts(fileExts)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the fileexts")
			}
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
			}
			c.Data(http.StatusOK, "application/json", gnuplex.MPV.SetPos(pos))
		}
	})
	gnuplex.Router.GET("/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, gnuplex.Last25())
	})
	gnuplex.Router.GET("/api/mediaitems", func(c *gin.Context) {
		res, err := gnuplex.NewDB.GetMediaItems()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/scanlib", func(c *gin.Context) {
		err := gnuplex.ScanLib()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})

}
