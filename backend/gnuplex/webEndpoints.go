package gnuplex

import (
	"gnuplex/consts"
	"gnuplex/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayMediaBody struct {
	Id       models.MediaItemId `json:"id"`
	PlayNext bool               `json:"play_next"`
	PlayLast bool               `json:"play_last"`
}

type CastBody struct {
	Url      string `json:"url"`
	Temp     bool   `json:"temp"`
	PlayNext bool   `json:"play_next"`
	PlayLast bool   `json:"play_last"`
}

type VolBody struct {
	Vol int `json:"vol"`
}

type PosBody struct {
	Pos int `json:"pos"`
}

type MediaItemsRes struct {
	Res   []models.MediaItem `json:"res"`
	Count int64              `json:"count"`
}

type MediaDirsBody []string

type FileExtsBody []string

type SubBody struct {
	Visible bool  `json:"visible"`
	ID      int64 `json:"id,omitempty"`
}

type FilterBody struct {
	Filter string `json:"filter"`
}

// Initialize the web server's HTTP Endpoints
func (gnuplex *GNUPlex) InitWebEndpoints(prod bool, staticFiles, sourceHash, platform, goVersion, exe string) {
	gnuplex.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	gnuplex.Router.GET("/gnuplex", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	gnuplex.Router.Static("/home", staticFiles)
	gnuplex.Router.GET("/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, consts.VersionInfo{Version: consts.Version, SourceHash: sourceHash, Platform: platform, GoVersion: goVersion})
	})
	gnuplex.Router.POST("/api/play", func(c *gin.Context) {
		if err := gnuplex.MPV.Play(); err != nil {
			log.Println("Error 36007463-94e1-4966-b718-8fbb7dd9de7c: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.POST("/api/pause", func(c *gin.Context) {
		if err := gnuplex.MPV.Pause(); err != nil {
			log.Println("Error 6dfeac1e-3b02-4667-8b8d-45d976857a22: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.GET("/api/nowplaying", func(c *gin.Context) {
		media, err := gnuplex.GetNowPlaying()
		if err != nil {
			log.Println("Error 6953e237-eec5-4701-aa41-1dbc20d0e4d2: ,", err)
			c.JSON(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, media)
		}
	})
	gnuplex.Router.POST("/api/playmedia", func(c *gin.Context) {
		body := PlayMediaBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error 2cd6e6a1-a62e-4336-83bb-b41c3e95f2ce: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.PlayById(body.Id, body.PlayNext, body.PlayLast); err != nil {
			log.Println("Error 2423754f-7a36-467a-86de-7fc05fb7a9b2: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.POST("/api/cast", func(c *gin.Context) {
		body := CastBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error 0bf4554b-ac42-487f-9468-fee0b1e002bd: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.Cast(body.Url, body.Temp, body.PlayNext, body.PlayLast); err != nil {
			log.Println("Error 65365414-ca17-4885-a50f-8584dd758512: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.GET("/api/paused", func(c *gin.Context) {
		paused, err := gnuplex.MPV.GetPaused()
		if err != nil {
			log.Println("Error 64ea9c3a-81d9-4102-ae82-5edd137e21e5: ,", err)
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, paused)
		}
	})
	gnuplex.Router.GET("/api/vol", func(c *gin.Context) {
		vol, err := gnuplex.MPV.GetVol()
		if err != nil {
			log.Println("Error f68be2e6-e6da-4424-aa25-15ba3d185a5a: ,", err)
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, vol)
		}
	})
	gnuplex.Router.POST("/api/vol", func(c *gin.Context) {
		body := VolBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error b5178076-4524-40a1-9523-04720ac09efd: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.MPV.SetVol(body.Vol); err != nil {
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.JSON(http.StatusOK, true)
		}
	})
	gnuplex.Router.GET("/api/mediadirs", func(c *gin.Context) {
		res, err := gnuplex.DB.GetMediaDirs()
		if err != nil {
			log.Println("Error 03224d1a-5df0-4ed3-a884-c494f7b23fea: ,", err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/mediadirs", func(c *gin.Context) {
		body := MediaDirsBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error 288fd456-7867-40bc-b880-3131e0655a4a: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.DB.SetMediadirs(body); err != nil {
			log.Println("Error 35e487a4-e426-453d-a668-24b592da655d: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.GET("/api/file_exts", func(c *gin.Context) {
		res, err := gnuplex.DB.GetFileExts()
		if err != nil {
			log.Println("Error a716951f-506b-41d4-bbd1-752316b1a827: ,", err)
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
			log.Println("Error 82a1db81-efec-4d17-9e27-c68094e8f239: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})

	gnuplex.Router.GET("/api/pos", func(c *gin.Context) {
		pos, err := gnuplex.MPV.GetPos()
		if err != nil {
			log.Println("Error c8065bdb-bf8b-466f-9fe3-bfeb70339390: ,", err)
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, pos)
		}
	})
	gnuplex.Router.GET("/api/timeremaining", func(c *gin.Context) {
		timeRemaining, err := gnuplex.MPV.GetTimeRemaining()
		if err != nil {
			log.Println("Error d2800412-f6b7-42e9-bc69-4de028a98252: ,", err)
			c.JSON(http.StatusInternalServerError, 0)
		} else {
			c.JSON(http.StatusOK, timeRemaining)
		}
	})
	gnuplex.Router.POST("/api/pos", func(c *gin.Context) {
		body := PosBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error 9296af2e-ab22-43b7-844b-795c593a4a4c: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
		} else if err = gnuplex.MPV.SetPos(body.Pos); err != nil {
			log.Println("Error 4dfb2f0c-36db-4aea-a156-737e9269140a: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}

	})
	gnuplex.Router.GET("/api/last25", func(c *gin.Context) {
		res, err := gnuplex.DB.GetLast25Played()
		if err != nil {
			log.Println("Error 3380b5ac-429a-4809-a660-70f38287ad7e: ,", err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}

	})
	gnuplex.Router.GET("/api/mediaitems", func(c *gin.Context) {
		search := c.Query("search")
		offsetStr := c.Query("offset")
		var offset int
		var err error
		if len(offsetStr) != 0 {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				log.Println("Error 59444453-f5aa-4f66-a82d-7ddee4027bed: ,", err)
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		res, count, err := gnuplex.DB.GetMediaItems(search, offset)
		if err != nil {
			log.Println("Error 1a56e258-f539-443e-99fc-74b24746d31e: ,", err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, MediaItemsRes{Res: res, Count: count})
		}
	})
	gnuplex.Router.POST("/api/scanlib", func(c *gin.Context) {
		if err := gnuplex.ScanLib(); err != nil {
			log.Println("Error 49374385-17c5-4348-a2c2-d6566fe53692: ", err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.GET("/api/sub", func(c *gin.Context) {
		res, err := gnuplex.GetSubs()
		if err != nil {
			log.Println("Error 33e3ed57-9da6-4740-8a81-4d8d56086af7: ,", err)
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
	gnuplex.Router.POST("/api/sub", func(c *gin.Context) {
		body := SubBody{}
		var err error
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error f49f4f29-753b-4b14-813f-59a38584fb2f: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
			return
		}
		if body.ID == -1 {
			err = gnuplex.SetSubVisibility(body.Visible)
		} else {
			if err = gnuplex.SetSubTrack(body.ID); err != nil {
				log.Println("Error 0ba2e0b1-92fb-47ce-aea7-806ae49de10b: ,", err)
				c.String(http.StatusInternalServerError, "some problem doing that")
				return
			}
			err = gnuplex.SetSubVisibility(true)
		}
		if err != nil {
			log.Println("Error fa1eb88e-1642-4c09-9f02-1ca49432a8d5: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else {
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.POST("/api/upgrade", func(c *gin.Context) {
		if upgraded, err := UpgradeGNUPlex(exe, false); err != nil {
			log.Println("Error 993a427d-e32d-48bc-8387-fe93840671f0: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
		} else if upgraded {
			log.Println("Upgraded and going into shutdown")
			c.Status(http.StatusOK)
			gnuplex.Wg.Done()
		} else {
			log.Println("Already up to date")
			c.Status(http.StatusOK)
		}
	})
	gnuplex.Router.POST("/api/filter", func(c *gin.Context) {
		body := FilterBody{}
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			log.Println("Error 0d997332-2815-43b4-9415-689e9a681cb8: ,", err)
			c.String(http.StatusBadRequest, "bad body format")
			return
		}
		if err := gnuplex.MPV.SetFilter(body.Filter); err != nil {
			log.Println("Error 20e52d4e-a82c-4d42-b4b8-16dcde63daf4: ,", err)
			c.String(http.StatusInternalServerError, "some problem doing that")
			return
		}
		c.Status(http.StatusOK)
	})
}
