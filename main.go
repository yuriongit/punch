package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// enable for production; disabled for debugging
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if err := router.Run(); err != nil {
		panic(err)
	}
}
