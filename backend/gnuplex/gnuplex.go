package gnuplex

import (
	"log"
	"sync"

	"gnuplex/db"
	"gnuplex/models"
	"gnuplex/mpv"

	"github.com/gin-gonic/gin"
)

type GNUPlex struct {
	DB        *db.DB
	Router    *gin.Engine
	PlayQueue [](*models.MediaItem)
	MPV       *mpv.MPV
}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket, dbPath, staticFiles string) (*GNUPlex, error) {
	/*
	 * HTTP backend
	 */
	gnuplex := new(GNUPlex)
	gnuplex.Router = gin.Default()
	gnuplex.Router.SetTrustedProxies(nil)
	gnuplex.InitWebEndpoints(verbose, staticFiles)
	/*
	 * mpv unix socket
	 */
	mpv, err := mpv.Init(wg, verbose, createMpvDaemon, mpvSocket)
	if err != nil {
		return nil, err
	}
	gnuplex.MPV = mpv
	/*
	 * new sqlite DB
	 */
	db, err := db.Init(dbPath, verbose)
	if err != nil {
		return nil, err
	}
	gnuplex.DB = db
	return gnuplex, nil
}

func (server *GNUPlex) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := server.Router.Run(":40000")
	if err != nil {
		log.Println("Server error:", err)
	}
	return err
}
