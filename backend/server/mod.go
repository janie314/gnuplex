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

func (server *Server) initEndpoints(api_url_base string) {
	server.Router.GET(api_url_base+"/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	server.Router.POST(api_url_base+"/play", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Play())
	})
	server.Router.POST(api_url_base+"/pause", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.Pause())
	})
	server.Router.POST(api_url_base+"/toggle", func(c *gin.Context) {
		paused, err := mpv.Toggle()
		if err != nil {
			c.JSON(http.StatusOK, paused)
		} else {
			c.JSON(http.StatusInternalServerError, nil)
		}
	})
	server.Router.GET(api_url_base+"/paused", func(c *gin.Context) {
		paused, err := mpv.IsPaused()
		if err != nil {
			c.JSON(http.StatusOK, paused)
		} else {
			c.JSON(http.StatusInternalServerError, nil)
		}
	})
	server.Router.GET(api_url_base+"/media", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetMedia())
	})
	server.Router.POST(api_url_base+"/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			server.AddHist(mediafile)
			c.Data(http.StatusOK, "application/json", mpv.SetMedia(mediafile))
		}
	})
	server.Router.GET(api_url_base+"/vol", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetVolume())
	})
	server.Router.POST(api_url_base+"/vol", func(c *gin.Context) {
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
	server.Router.GET(api_url_base+"/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMediadirs(false))
	})
	server.Router.POST(api_url_base+"/mediadirs", func(c *gin.Context) {
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
	server.Router.GET(api_url_base+"/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetFileExts(false))
	})
	server.Router.POST(api_url_base+"/file_exts", func(c *gin.Context) {
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
	server.Router.GET(api_url_base+"/pos", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", mpv.GetPos())
	})
	server.Router.POST(api_url_base+"/pos", func(c *gin.Context) {
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
	server.Router.GET(api_url_base+"/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.Last25())
	})
	server.Router.GET(api_url_base+"/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, server.GetMedialib(false))
	})
	server.Router.POST(api_url_base+"/medialist", func(c *gin.Context) {
		server.ScanLib(false)
		c.String(http.StatusOK, "OK")
	})
}

func New(wg *sync.WaitGroup, prod bool, port int, static_url_base, api_url_base, db_path, static_dir_path string) (*Server, error) {
	server := new(Server)
	server.Router = gin.Default()
	server.Router.SetTrustedProxies(nil)
	go mpv.InitUnixConn(wg)
	db, err := liteDB.New(prod, db_path)
	if err != nil {
		return nil, err
	}
	server.DB = db
	server.port = port
	/*
	 * Serve static files
	 */
	server.Router.Static(static_url_base, static_dir_path)
	/*
	 * API endpoints
	 */
	server.initEndpoints(api_url_base)
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
