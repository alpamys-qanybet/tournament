package api

import (
	"net/http"

	"tournament/internal/ctrl"

	"github.com/gin-gonic/gin"
)

func GetDivisions(c *gin.Context) {
	res, err := ctrl.GetDivisions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func PrepareDivisions(c *gin.Context) {

	err := ctrl.PrepareDivisions(c.Request.Context())
	if err != nil {

		errMsg := err.Error()

		if errMsg == "must_have_16_teams_to_prepare_divisions" || errMsg == "division_is_already_prepared" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"err": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func StartDivisions(c *gin.Context) {

	err := ctrl.StartDivisions(c.Request.Context())
	if err != nil {

		errMsg := err.Error()

		if errMsg == "must_have_16_teams_to_start_divisions" || errMsg == "division_is_not_prepared" || errMsg == "division_is_already_started" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"err": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
