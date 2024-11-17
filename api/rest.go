package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func extractBody(c *gin.Context, data *map[string]interface{}) error {
	body := c.Request.Body
	value, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(value), &data)
	return nil
}

func RootIndex(c *gin.Context) {
	c.String(http.StatusOK, "it works")
}
