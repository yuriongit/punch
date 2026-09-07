package loadtest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuriongit/punch/shared/types"
)

// RegisterTest creates a test and persists it to Redis.
func RegisterTest(c *gin.Context) {
	var config types.PunchConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Println(config)

	testID, err := createTest(&config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create test",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"test_id": testID,
	})
}
