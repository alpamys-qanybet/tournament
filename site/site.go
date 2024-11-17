package site

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppHome(c *gin.Context) {
	c.HTML(http.StatusOK, "app", nil)
}
