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
	"gnuplex-backend/models"
	"gnuplex-backend/mpv"
	"gnuplex-backend/server/liteDB"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	DB        *liteDB.LiteDB
	NewDB     *gorm.DB
	Router    *gin.Engine
	PlayQueue [](*models.MediaItem)
	MPV       *net.UnixConn
	MPVMutex  *sync.Mutex
}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket, dbPath string) (*Server, error) {
	/*
	 * HTTP backend
	 */
	server := new(Server)
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
func (server *Server) initEndpoints(prod bool) {
	server.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	server.Router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	if prod {
		server.Router.Static("/home", consts.ProdStaticFilespath)
	} else {
		server.Router.Static("/home", consts.DevStaticFilespath)
	}
	server.Router.GET("/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	server.Router.POST("/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Play(server.MPV, server.MPVMutex))
	})
	server.Router.POST("/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Pause(server.MPV, server.MPVMutex))
	})
	server.Router.GET("/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.IsPaused(server.MPV, server.MPVMutex))
	})
	server.Router.GET("/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetMedia(server.MPV, server.MPVMutex))
	})
	server.Router.POST("/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			server.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", mpv.SetMedia(server.MPV, server.MPVMutex, mediafile))
		}
	})
	server.Router.GET("/api/vol", func(c *gin.Context) {
		vol, err := mpv.GetVol(server.MPV, server.MPVMutex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})
	server.Router.POST("/api/vol", func(c *gin.Context) {
		param := c.Query("vol")
		if param == "" {
			c.String(http.StatusBadRequest, "empty vol string")
		} else {
			vol, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad vol string")
			}
			c.Data(http.StatusOK, "application/json", mpv.SetVolume(server.MPV, server.MPVMutex, vol))
		}
	})
	server.Router.GET("/api/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMediadirs(false))
	})
	server.Router.POST("/api/mediadirs", func(c *gin.Context) {
		mediadirsJson := []byte(c.Query("mediadirs"))
		var mediadirs []string
		err := json.Unmarshal(mediadirsJson, &mediadirs)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = server.SetMediadirs(mediadirs)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the mediadirs")
			}
		}
	})
	server.Router.GET("/api/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetFileExts(false))
	})
	server.Router.POST("/api/file_exts", func(c *gin.Context) {
		fileExtsJson := []byte(c.Query("file_exts"))
		var fileExts []string
		err := json.Unmarshal(fileExtsJson, &fileExts)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = server.SetFileExts(fileExts)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the fileexts")
			}
		}
	})
	server.Router.GET("/api/pos", func(c *gin.Context) {
		vol, err := mpv.GetPos(server.MPV, server.MPVMutex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})

	server.Router.GET("/api/timeremaining", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetTimeRemaining(server.MPV, server.MPVMutex))
	})
	server.Router.POST("/api/pos", func(c *gin.Context) {
		param := c.Query("pos")
		if param == "" {
			c.String(http.StatusBadRequest, "empty pos string")
		} else {
			pos, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad pos string")
			}
			c.Data(http.StatusOK, "application/json", mpv.SetPos(server.MPV, server.MPVMutex, pos))
		}
	})
	server.Router.GET("/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.Last25())
	})
	server.Router.GET("/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMedialib(false))
	})
	server.Router.POST("/api/medialist", func(c *gin.Context) {
		server.ScanLib()
		c.String(http.StatusOK, "OK")
	})

}

func (server *Server) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := server.Router.Run(":40000")
	if err != nil {
		log.Println("Server error:", err)
	}
	return err
}
