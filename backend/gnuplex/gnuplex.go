package gnuplex

import (
	"fmt"
	"log"
	"sync"

	"gnuplex/db"
	"gnuplex/models"
	"gnuplex/mpv"

	"github.com/gin-gonic/gin"
)

type GNUPlex struct {
	DB        *db.DB
	Port      int
	Router    *gin.Engine
	PlayQueue [](*models.MediaItem)
	MPV       *mpv.MPV
}

// Initialize a GNUPlex instance.
func Init(wg *sync.WaitGroup, verbose bool, dbPath, staticFiles string, port int, sourceHash string) (*GNUPlex, error) {
	// HTTP backend
	gnuplex := new(GNUPlex)
	gnuplex.Router = gin.Default()
	gnuplex.Router.SetTrustedProxies(nil)
	gnuplex.InitWebEndpoints(verbose, staticFiles, sourceHash)
	// MPV instance
	mpv, err := mpv.Init(wg, verbose)
	if err != nil {
		return nil, err
	}
	gnuplex.MPV = mpv
	// SQLite DB
	db, err := db.Init(dbPath, verbose)
	if err != nil {
		return nil, err
	}
	gnuplex.DB = db
	gnuplex.Port = port
	return gnuplex, nil
}

// Run the GNUPlex daemon.
func (server *GNUPlex) Run(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := server.Router.Run(fmt.Sprintf(":%d", server.Port))
	if err != nil {
		log.Println("Server error:", err)
	}
	return err
}
