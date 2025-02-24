package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tgnotifyhub/config"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader(config.Loaded().AuthHeader)

	if token != config.Loaded().AuthToken {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
