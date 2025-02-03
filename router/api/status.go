package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/CoiaPrant/Sleepy/services/device"
)

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Ok":   true,
		"Data": device.ListDevices(),
	})
}
