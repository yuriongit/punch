// Package main implements the functionality for the API Punch provides.
// This includes an orchestrator and workers to complete the distributed HTTP load test.
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yuriongit/punch/internal/infra"
	"github.com/yuriongit/punch/internal/services/health"
	"github.com/yuriongit/punch/internal/services/loadtest"
)

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
