package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"gnuplex-backend/consts"
	"gnuplex-backend/db"
	"gnuplex-backend/liteDB"
	"gnuplex-backend/models"
	"gnuplex-backend/mpv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GNUPlex struct {
	DB        *liteDB.LiteDB
	NewDB     *gorm.DB
	Router    *gin.Engine
	PlayQueue [](*models.MediaItem)
	MPV       *net.UnixConn
	MPVMutex  *sync.Mutex
}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket, dbPath string) (*GNUPlex, error) {
	/*
	 * HTTP backend
	 */
	server := new(GNUPlex)
	server.Router = gin.Default()
	server.Router.SetTrustedProxies(nil)
	server.initEndpoints(verbose)
	/*
	 * mpv unix socket
	 */
	var mpvConn *net.UnixConn
	mpvConn, err := mpv.Init(wg, verbose, createMpvDaemon, mpvSocket)
	if err != nil {
		return nil, err
	}
	server.MPV = mpvConn
	mu := &sync.Mutex{}
	server.MPVMutex = mu
	/*
	 * old sqlite DB
	 */
	oldDb, err := liteDB.Init(verbose)
	if err != nil {
		return nil, err
	}
	server.DB = oldDb
	/*
	 * new sqlite DB
	 */
	newDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	server.NewDB = newDB
	if err = db.Init(server.NewDB); err != nil {
		return nil, err
	}
	return server, nil
}

// Initialize the web server's HTTP Endpoints
func (gnuplex *GNUPlex) initEndpoints(prod bool) {
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
		c.Data(http.StatusOK, "application/json", mpv.Play(gnuplex.MPV, gnuplex.MPVMutex))
	})
	gnuplex.Router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Pause(gnuplex.MPV, gnuplex.MPVMutex))
	})
	gnuplex.Router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.IsPaused(gnuplex.MPV, gnuplex.MPVMutex))
	})
	gnuplex.Router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetMedia(gnuplex.MPV, gnuplex.MPVMutex))
	})
	gnuplex.Router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			gnuplex.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", mpv.SetMedia(gnuplex.MPV, gnuplex.MPVMutex, mediafile))
		}
	})
	gnuplex.Router.GET("/api/vol", func(c *gin.Context) {
		vol, err := mpv.GetVol(gnuplex.MPV, gnuplex.MPVMutex)
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
			c.Data(http.StatusOK, "application/json", mpv.SetVolume(gnuplex.MPV, gnuplex.MPVMutex, vol))
		}
	})
	gnuplex.Router.GET("/api/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gnuplex.GetMediadirs(false))
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
		vol, err := mpv.GetPos(gnuplex.MPV, gnuplex.MPVMutex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})

	gnuplex.Router.GET("/api/timeremaining", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetTimeRemaining(gnuplex.MPV, gnuplex.MPVMutex))
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
			c.Data(http.StatusOK, "application/json", mpv.SetPos(gnuplex.MPV, gnuplex.MPVMutex, pos))
		}
	})
	gnuplex.Router.GET("/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, gnuplex.Last25())
	})
	gnuplex.Router.GET("/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, gnuplex.GetMedialib(false))
	})
	gnuplex.Router.POST("/api/medialist", func(c *gin.Context) {
		gnuplex.ScanLib()
		c.String(http.StatusOK, "OK")
	})

}

func (server *GNUPlex) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := server.Router.Run(":40000")
	if err != nil {
		log.Println("Server error:", err)
	}
	return err
}
