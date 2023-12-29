package server

import (
	"encoding/json"
	"gnuplex-backend/consts"
	"gnuplex-backend/mpv"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) initEndpoints(api_url_base string) {
	server.Router.GET(api_url_base+"/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	server.Router.GET(api_url_base+"/paused", func(c *gin.Context) {
		paused, err := server.mpv.IsPaused()
		readQuery2HTTP(c, paused, err)
	})
	server.Router.POST(api_url_base+"/toggle", func(c *gin.Context) {
		paused, err := server.mpv.Toggle()
		readQuery2HTTP(c, paused, err)
	})
	server.Router.GET(api_url_base+"/media", func(c *gin.Context) {
		media, err := server.mpv.GetMedia()
		readQuery2HTTP(c, media, err)
	})
	server.Router.POST(api_url_base+"/media", func(c *gin.Context) {
		mediafile := c.Query("mediafile")
		if mediafile == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			err := server.mpv.SetMedia(mediafile)
			if err == nil {
				server.DB.AddHist(mediafile, false)
			}
			writeQuery2HTTP(c, err)
		}
	})
	server.Router.GET(api_url_base+"/vol", func(c *gin.Context) {
		vol, err := server.mpv.GetVolume()
		readQuery2HTTP(c, vol, err)
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
			writeQuery2HTTP(c, server.mpv.SetVolume(float64(vol)))
		}
	})
	// TODO fold this into /pos
	server.Router.POST(api_url_base+"/incpos", func(c *gin.Context) {
		param := c.Query("inc")
		if param == "" {
			c.String(http.StatusBadRequest, "empty inc string")
		} else {
			inc, err := strconv.Atoi(param)
			if err != nil {
				c.String(http.StatusBadRequest, "bad inc string")
			}
			writeQuery2HTTP(c, server.mpv.IncPos(float64(inc)))
		}
	})
	// TODO fold into new framework
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
		pos, err := server.mpv.GetPos()
		readQuery2HTTP(c, pos, err)
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
			err = server.mpv.SetPos(float64(pos))
			if err != nil {
				c.Data(http.StatusOK, "application/json", nil)
			} else {
				c.Data(http.StatusOK, "application/json", nil)
			}
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

func readQuery2HTTP[T mpv.ResponseData](c *gin.Context, val T, err error) {
	if err != nil {
		log.Println("Error", err)
		c.JSON(http.StatusInternalServerError, nil)
	} else {
		c.JSON(http.StatusOK, val)
	}
}

func writeQuery2HTTP(c *gin.Context, err error) {
	if err != nil {
		log.Println("Error", err)
		c.JSON(http.StatusInternalServerError, nil)
	} else {
		c.JSON(http.StatusOK, nil)
	}
}
