package core

import "github.com/gin-gonic/gin"

func Check(err error, c *gin.Context, status int) bool {
	if err == nil {
		return true
	}

	c.AbortWithStatus(status)
	return false
}
