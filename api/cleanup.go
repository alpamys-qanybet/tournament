package api

import (
	"net/http"

	"tournament/internal/ctrl"

	"github.com/gin-gonic/gin"
)

func Cleanup(c *gin.Context) {

	err := ctrl.Cleanup(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
