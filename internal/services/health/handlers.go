// Package health provides health check handlers for the server.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Live checks whether the server is alive.
func Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"code": "HEALTHY",
	})
}
