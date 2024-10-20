package gnuplex

import (
	"log"
	"sync"

	"gnuplex-backend/db"
	"gnuplex-backend/liteDB"
	"gnuplex-backend/models"
	"gnuplex-backend/mpv"

	"github.com/gin-gonic/gin"
)

type GNUPlex struct {
	DB        *liteDB.LiteDB
	NewDB     *db.DB
	Router    *gin.Engine
	PlayQueue [](*models.MediaItem)
	MPV       *mpv.MPV
}

func Init(wg *sync.WaitGroup, verbose, createMpvDaemon bool, mpvSocket, dbPath string) (*GNUPlex, error) {
	/*
	 * HTTP backend
	 */
	gnuplex := new(GNUPlex)
	gnuplex.Router = gin.Default()
	gnuplex.Router.SetTrustedProxies(nil)
	gnuplex.InitWebEndpoints(verbose)
	/*
	 * mpv unix socket
	 */
	mpv, err := mpv.Init(wg, verbose, createMpvDaemon, mpvSocket)
	if err != nil {
		return nil, err
	}
	gnuplex.MPV = mpv
	/*
	 * old sqlite DB
	 */
	oldDb, err := liteDB.Init(verbose)
	if err != nil {
		return nil, err
	}
	gnuplex.DB = oldDb
	/*
	 * new sqlite DB
	 */
	newDB, err := db.Init(dbPath)
	if err != nil {
		return nil, err
	}
	gnuplex.NewDB = newDB
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
