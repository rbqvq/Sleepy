package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Check(c *gin.Context) {
	if !c.IsWebsocket() {
		c.AbortWithStatus(http.StatusUpgradeRequired)
		return
	}
}

func NoRoute(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func NoMethod(c *gin.Context) {
	c.Status(http.StatusMethodNotAllowed)
}
