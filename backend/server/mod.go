package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"gnuplex-backend/consts"
	"gnuplex-backend/liteDB"
	"gnuplex-backend/mpv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB     *liteDB.LiteDB
	Router *gin.Engine
	port   int
}

func (server *Server) initEndpoints(url_base string) {
	server.Router.GET(url_base+"/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	server.Router.POST(url_base+"/api/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Play())
	})
	server.Router.POST(url_base+"/api/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Pause())
	})
	server.Router.GET(url_base+"/api/paused", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.IsPaused())
	})
	server.Router.GET(url_base+"/api/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetMedia())
	})
	server.Router.POST(url_base+"/api/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			server.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", mpv.SetMedia(mediafile))
		}
	})
	server.Router.GET(url_base+"/api/vol", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetVolume())
	})
	server.Router.POST(url_base+"/api/vol", func(c *gin.Context) {
		param := c.Query("vol")
		if param == "" {
			c.String(http.StatusBadRequest, "empty vol string")
		} else {
			vol, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad vol string")
			}
			c.Data(http.StatusOK, "application/json", mpv.SetVolume(vol))
		}
	})
	server.Router.GET(url_base+"/api/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMediadirs(false))
	})
	server.Router.POST(url_base+"/api/mediadirs", func(c *gin.Context) {
		mediadirsJson := []byte(c.Query("mediadirs"))
		var mediadirs []string
		err := json.Unmarshal(mediadirsJson, &mediadirs)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = server.SetMediadirs(mediadirs, false)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the mediadirs")
			}
		}
	})
	server.Router.GET(url_base+"/api/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetFileExts(false))
	})
	server.Router.POST(url_base+"/api/file_exts", func(c *gin.Context) {
		fileExtsJson := []byte(c.Query("file_exts"))
		var fileExts []string
		err := json.Unmarshal(fileExtsJson, &fileExts)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = server.SetFileExts(fileExts, false)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the fileexts")
			}
		}
	})
	server.Router.GET(url_base+"/api/pos", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetPos())
	})
	server.Router.POST(url_base+"/api/pos", func(c *gin.Context) {
		param := c.Query("pos")
		if param == "" {
			c.String(http.StatusBadRequest, "empty pos string")
		} else {
			pos, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad pos string")
			}
			c.Data(http.StatusOK, "application/json", mpv.SetPos(pos))
		}
	})
	server.Router.GET(url_base+"/api/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.Last25())
	})
	server.Router.GET(url_base+"/api/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMedialib(false))
	})
	server.Router.POST(url_base+"/api/medialist", func(c *gin.Context) {
		server.ScanLib(false)
		c.String(http.StatusOK, "OK")
	})
}

func New(wg *sync.WaitGroup, prod bool, port int, base_url, path string) (*Server, error) {
	server := new(Server)
	server.Router = gin.Default()
	server.Router.SetTrustedProxies(nil)
	go mpv.InitUnixConn(wg)
	db, err := liteDB.New(prod, path)
	if err != nil {
		return nil, err
	}
	server.DB = db
	server.port = port
	/*
	 * Serve static files
	 */
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
	/*
	 * API endpoints
	 */
	server.initEndpoints(base_url)
	return server, nil
}

func (server *Server) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := server.Router.Run(fmt.Sprintf(":%d", server.port))
	if err != nil {
		log.Println("Server error:", err)
	}
	return err
}
