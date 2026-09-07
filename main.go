// Package main implements the functionality for the API Punch provides.
// This includes an orchestrator and workers to complete the distributed HTTP load test.
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuriongit/punch/infra"
	"github.com/yuriongit/punch/internal/services/loadtest"
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
	r := gin.Default()
	if err := infra.InitRedis(); err != nil {
		panic(err)
	}

	// Health routes group
	{
		r.Group("/health")
		r.GET("/live", HealthLive)
	}
	{
	  lt := r.Group("/load-test")
		lt.POST("/register", loadtest.RegisterTest)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
