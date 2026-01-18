package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteBindError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": gin.H{
			"code":    "BAD_REQUEST",
			"message": "invalid request payload",
			"details": gin.H{
				"cause": err.Error(),
			},
		},
	})
}
