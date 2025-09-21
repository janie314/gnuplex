package gnuplex

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	Wg        *sync.WaitGroup
}

// Initialize a GNUPlex instance.
func Init(wg *sync.WaitGroup, verbose bool, dbPath, staticFiles string, port int, sourceHash, platform, exe string) (*GNUPlex, error) {
	// HTTP backend
	gnuplex := new(GNUPlex)
	gnuplex.Router = gin.Default()
	gnuplex.Router.SetTrustedProxies(nil)
	gnuplex.InitWebEndpoints(verbose, staticFiles, sourceHash, platform, exe)
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
	gnuplex.Wg = wg
	return gnuplex, nil
}

// Run the GNUPlex daemon.
func (server *GNUPlex) Run(wg *sync.WaitGroup) error {
	defer server.Wg.Done()
	err := server.Router.Run(fmt.Sprintf(":%d", server.Port))
	if err != nil {
		log.Println("e6c1f9ee-681d-4eb4-b487-d86622e07aac Server error:", err)
	}
	return err
}

// Upgrade GNUPlex
func UpgradeGNUPlex(exe string, interactive bool) error {
	cmd := exec.Command("git", "-C", filepath.Join(filepath.Dir(exe), "../.."), "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	} else if interactive {
		fmt.Println("Successfully upgraded! Now run `systemctl --user restart gnuplex`.")
	}
	return nil
}
