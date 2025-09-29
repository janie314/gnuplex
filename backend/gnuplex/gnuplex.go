package gnuplex

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
func Init(wg *sync.WaitGroup, verbose bool, dbPath, staticFiles string, port int, sourceHash, platform, goVersion, exe string) (*GNUPlex, error) {
	// HTTP backend
	gnuplex := new(GNUPlex)
	gnuplex.Router = gin.Default()
	gnuplex.Router.SetTrustedProxies(nil)
	gnuplex.InitWebEndpoints(verbose, staticFiles, sourceHash, platform, goVersion, exe)
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
//
// The boolean in the return code represents whether we upgraded or not
func UpgradeGNUPlex(exe string, interactive bool) (bool, error) {
	cwd := filepath.Join(filepath.Dir(exe), "../..")
	branchB, err := exec.Command("git", "-C", cwd, "branch", "--show-current").Output()
	if err != nil {
		return false, fmt.Errorf("19c5e798-666b-4aae-818b-56061645ca9f git problem %v", err)
	}
	branch := strings.TrimSpace(string(branchB))
	if err := exec.Command("git", "-C", cwd, "fetch", "origin", branch).Run(); err != nil {
		return false, fmt.Errorf("fetch failed %v", err)
	}
	local, err := exec.Command("git", "-C", cwd, "rev-parse", branch).Output()
	if err != nil {
		return false, err
	}
	remote, err := exec.Command("git", "-C", cwd, "rev-parse", "origin/"+branch).Output()
	if err != nil {
		return false, err
	}
	if bytes.Equal(local, remote) {
		if interactive {
			fmt.Println("Nothing to upgrade")
		}
		return false, nil
	}
	cmd := exec.Command("git", "-C", cwd, "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return false, err
	} else if interactive {
		fmt.Println("Successfully upgraded! Now run `systemctl --user restart gnuplex`.")
	}
	return true, nil
}
