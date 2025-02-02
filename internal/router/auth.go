package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	AccountId uint64 `json:"account_id" binding:"required"`
}

func HandleAuth(c *gin.Context) {
	var request AuthRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
