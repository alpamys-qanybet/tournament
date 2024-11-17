package api

import (
	"net/http"

	"tournament/internal/ctrl"

	"github.com/gin-gonic/gin"
)

func GetTeamList(c *gin.Context) {
	res, err := ctrl.GetTeamList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateTeam(c *gin.Context) {

	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "invalid_body_params",
		})
		return
	}

	var name string
	if bodyData["name"] != nil {
		name = bodyData["name"].(string)
	}

	id, err := ctrl.CreateTeam(c.Request.Context(), name)
	if err != nil {
		errMsg := err.Error()

		if errMsg == "create_team_failure_name_is_required" || errMsg == "max_16_teams_allowed" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"err": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errMsg,
		})
		return
	}

	data := gin.H{
		"id": id,
	}

	c.JSON(http.StatusCreated, data)
}

func GenerateTeams(c *gin.Context) {

	err := ctrl.GenerateTeams(c.Request.Context())
	if err != nil {
		errMsg := err.Error()

		if errMsg == "generation_only_allowed_into_empty_table" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"err": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
