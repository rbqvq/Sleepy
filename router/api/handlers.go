package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"Ok":   false,
		"Msg": "404 not found",
	})
}

func NoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed,gin.H{
		"Ok":   false,
		"Msg": "405 method not allowed",
	})
}
