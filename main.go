// Package main implements the functionality for the API Punch provides.
// This includes an orchestrator and workers to complete the distributed HTTP load test.
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes

// HealthLive checks whether the server is alive.
func HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func main() {
	// enable for production; disabled for debugging
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if err := router.Run(); err != nil {
		panic(err)
	}
}
