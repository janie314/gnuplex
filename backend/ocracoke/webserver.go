package ocracoke

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"gnuplex-backend/consts"
	"gnuplex-backend/mpvdaemon/mpvcmd"
	"gnuplex-backend/ocracoke/liteDB"

	"github.com/gin-gonic/gin"
)

type Ocracoke struct {
	DB     *liteDB.LiteDB
	Router *gin.Engine
}

func Init(wg *sync.WaitGroup, prod bool, mpvSocket string) (*Ocracoke, error) {
	oc := new(Ocracoke)
	oc.Router = gin.Default()
	oc.Router.SetTrustedProxies(nil)
	go mpvcmd.InitUnixConn(wg, mpvSocket)
	db, err := liteDB.Init(prod)
	if err != nil {
		return nil, err
	}
	oc.DB = db
	/*
	 * Serve static files
	 */
	oc.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	oc.Router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	if prod {
		oc.Router.Static("/home", consts.ProdStaticFilespath)
	} else {
		oc.Router.Static("/home", consts.DevStaticFilespath)
	}
	/*
	 * API endpoints
	 */
	oc.Router.GET("/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	oc.Router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Play())
	})
	oc.Router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.Pause())
	})
	oc.Router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.IsPaused())
	})
	oc.Router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetMedia())
	})
	oc.Router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			oc.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", mpvcmd.SetMedia(mediafile))
		}
	})
	oc.Router.GET("/api/vol", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetVolume())
	})
	oc.Router.POST("/api/vol", func(c *gin.Context) {
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
	oc.Router.GET("/api/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, oc.GetMediadirs(false))
	})
	oc.Router.POST("/api/mediadirs", func(c *gin.Context) {
		mediadirsJson := []byte(c.Query("mediadirs"))
		var mediadirs []string
		err := json.Unmarshal(mediadirsJson, &mediadirs)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = oc.SetMediadirs(mediadirs)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the mediadirs")
			}
		}
	})
	oc.Router.GET("/api/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, oc.GetFileExts(false))
	})
	oc.Router.POST("/api/file_exts", func(c *gin.Context) {
		fileExtsJson := []byte(c.Query("file_exts"))
		var fileExts []string
		err := json.Unmarshal(fileExtsJson, &fileExts)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = oc.SetFileExts(fileExts)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the fileexts")
			}
		}
	})
	oc.Router.GET("/api/pos", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetPos())
	})
	oc.Router.GET("/api/timeremaining", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpvcmd.GetTimeRemaining())
	})
	oc.Router.POST("/api/pos", func(c *gin.Context) {
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
	oc.Router.GET("/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, oc.Last25())
	})
	oc.Router.GET("/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, oc.GetMedialib(false))
	})
	oc.Router.POST("/api/medialist", func(c *gin.Context) {
		oc.ScanLib()
		c.String(http.StatusOK, "OK")
	})
	return oc, nil
}

func (ocracoke *Ocracoke) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := ocracoke.Router.Run(":40000")
	if err != nil {
		log.Println("Ocracoke error:", err)
	}
	return err
}
