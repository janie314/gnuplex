package server

import (
	"encoding/json"
	"gnuplex-backend/consts"
	"gnuplex-backend/mpv"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Types for various request and response JSON bodies
 */
type mediaBody struct {
	Media string `json:"media"`
}

type volBody struct {
	Vol float64 `json:"vol"`
}

type posBody struct {
	Pos float64 `json:"pos"`
	Inc bool    `json:"inc"`
}

type posResponse struct {
	Pos    float64 `json:"pos"`
	MaxPos float64 `json:"max_pos"`
}

type mediadirsBody struct {
	Mediadirs []string `json:"mediadirs"`
}

type fileExtsBody struct {
	FileExts []string `json:"file_exts"`
}

func (srv *Server) initEndpoints(api_url_base string) {
	/*
	 * GET /version
	 * Response: string
	 *   GNUPlex version string.
	 */
	srv.Router.GET(api_url_base+"/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.GNUPlexVersion)
	})
	/*
	 * GET /paused
	 * Response: boolean
	 *   If current media is paused.
	 */
	srv.Router.GET(api_url_base+"/paused", func(c *gin.Context) {
		paused, err := srv.mpv.IsPaused()
		readQueryResponse(c, paused, err)
	})
	/*
	 * POST /paused
	 * Toggles video's play/pause status.
	 * Response: boolean
	 *   If current media is paused.
	 */
	srv.Router.POST(api_url_base+"/toggle", func(c *gin.Context) {
		paused, err := srv.mpv.Toggle()
		readQueryResponse(c, paused, err)
	})
	/*
	 * GET /media
	 * Response: string
	 *   Current media file.
	 */
	srv.Router.GET(api_url_base+"/media", func(c *gin.Context) {
		media, err := srv.mpv.GetMedia()
		readQueryResponse(c, media, err)
	})
	/*
	 * POST /media
	 * Body:
	 *   {
	 *     media: string;
	 *   }
	 */
	srv.Router.POST(api_url_base+"/media", func(c *gin.Context) {
		var media mediaBody
		c.BindJSON(&media)
		if media.Media == "" {
			c.String(http.StatusBadRequest, "empty mediafile string")
		} else {
			err := srv.mpv.SetMedia(media.Media)
			if err == nil {
				srv.DB.AddHist(media.Media, false)
			}
			writeQueryResponse(c, err)
		}
	})
	/*
	 * GET /vol
	 * Response: number
	 *   Current volume.
	 */
	srv.Router.GET(api_url_base+"/vol", func(c *gin.Context) {
		vol, err := srv.mpv.GetVolume()
		readQueryResponse(c, vol, err)
	})
	/*
	 * POST /vol
	 * Body:
	 *   {
	 *     vol: number;
	 *   }
	 */
	srv.Router.POST(api_url_base+"/vol", func(c *gin.Context) {
		var vol volBody
		c.BindJSON(&vol)
		writeQueryResponse(c, srv.mpv.SetVolume(vol.Vol))
	})
	/*
	* GET /pos
	* Response:
	*   {
	*     pos: number;     # both have units of seconds
	*     max_pos: number;
	*   }
	 */
	srv.Router.GET(api_url_base+"/pos", func(c *gin.Context) {
		var res posResponse
		pos, err := srv.mpv.GetPos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
		} else {
			res.Pos = pos
			remaining, err := srv.mpv.GetTimeRemaining()
			if err != nil {
				readQueryResponse(c, 0.0, err)
			} else {
				res.MaxPos = pos + remaining
				c.JSON(http.StatusOK, res)
			}
		}
	})
	/*
	 * POST /pos
	 * Body:
	 *   {
	 *     pos: number;
	 *     inc: boolean;   // whether this is an absolute position or an increment from the current position
	 *   }
	 * Response: if no error, responds with the following body with the current position:
	 *   {
	 *     pos: number;
	 *   }
	 */
	srv.Router.POST(api_url_base+"/pos", func(c *gin.Context) {
		var pos posBody
		err := c.BindJSON(&pos)
		if err != nil {
			log.Println("err", err)
			c.JSON(http.StatusBadRequest, nil)
		} else {
			if pos.Inc {
				err = srv.mpv.IncPos(pos.Pos)
			} else {
				err = srv.mpv.SetPos(pos.Pos)
			}
			if err != nil {
				newPos, err := srv.mpv.GetPos()
				readQueryResponse(c, newPos, err)
			} else {
				c.Data(http.StatusOK, "application/json", nil)
			}
		}
	})
	srv.Router.GET(api_url_base+"/last25", func(c *gin.Context) {
		c.JSON(http.StatusOK, srv.DB.Last25(false))
	})
	srv.Router.GET(api_url_base+"/medialist", func(c *gin.Context) {
		c.JSON(http.StatusOK, srv.DB.GetMedialib(false))
	})
	srv.Router.POST(api_url_base+"/medialist", func(c *gin.Context) {
		srv.DB.ScanLib(false)
		c.String(http.StatusOK, "OK")
	})
	srv.Router.GET(api_url_base+"/mediadirs", func(c *gin.Context) {
		c.JSON(http.StatusOK, srv.DB.GetMediadirs(false))
	})
	srv.Router.POST(api_url_base+"/mediadirs", func(c *gin.Context) {
		// TODO use this `body` naming convention
		var body mediadirsBody
		err := c.BindJSON(&body)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = srv.DB.SetMediadirs(body.Mediadirs, false)
			if err == nil {
				c.JSON(http.StatusOK, nil)
			} else {
				c.JSON(http.StatusInternalServerError, nil)
			}
		}
	})
	srv.Router.GET(api_url_base+"/file_exts", func(c *gin.Context) {
		c.JSON(http.StatusOK, srv.DB.GetFileExts(false))
	})
	srv.Router.POST(api_url_base+"/file_exts", func(c *gin.Context) {
		var body fileExtsBody
		err := c.BindJSON(&body)
		if err != nil {
			c.String(http.StatusBadRequest, "bad mediadirs string")
		} else {
			err = srv.DB.SetFileExts(body.FileExts, false)
			if err == nil {
				c.JSON(http.StatusOK, "ok")
			} else {
				c.JSON(http.StatusInternalServerError, "Couldn't add the fileexts")
			}
		}
	})

}

// cast an mpvcmd read query into a Gin response
func readQueryResponse[T mpv.ResponseDatum](c *gin.Context, val T, err error) {
	if err != nil {
		log.Println("Error", err)
		c.JSON(http.StatusInternalServerError, nil)
	} else {
		c.JSON(http.StatusOK, val)
	}
}

// cast an mpvcmd write query into a Gin response
func writeQueryResponse(c *gin.Context, err error) {
	if err != nil {
		log.Println("Error", err)
		c.JSON(http.StatusInternalServerError, nil)
	} else {
		c.JSON(http.StatusOK, nil)
	}
}
