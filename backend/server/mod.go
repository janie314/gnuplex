package server

import (
	"fmt"
	"gnuplex-backend/liteDB"
	"gnuplex-backend/mpv"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB     *liteDB.LiteDB
	Router *gin.Engine
	port   int
	mpv    *mpv.MPV
}

func New(wg *sync.WaitGroup, prod bool, port int, static_url_base, api_url_base, db_path, static_dir_path string) (*Server, error) {
	server := new(Server)
	server.Router = gin.Default()
	server.Router.SetTrustedProxies(nil)
	mpv.New(wg)
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
