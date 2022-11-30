package webserver

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func Run(wg *sync.WaitGroup) {
	defer wg.Done()
	router := gin.Default()
	/*
	 * API endpoints
	 */
	router.POST("/play", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	/*
	 * Serve static files
	 */
	router.Static("/", "./public")
	/*
	 * Main execution
	 */
	router.Run(":50000")
}
