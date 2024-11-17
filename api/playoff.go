package api

import (
	"net/http"
	"tournament/internal/ctrl"

	"github.com/gin-gonic/gin"
)

func GetPlayoffs(c *gin.Context) {
	res, err := ctrl.GetPlayoffs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func PreparePlayoff(c *gin.Context) {

	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "invalid_body_params",
		})
		return
	}

	var stage string
	if bodyData["stage"] != nil {
		stage = bodyData["stage"].(string)
	}

	if stage == "quarter" {
		err := ctrl.PreparePlayoffQuarter(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "division_is_not_started" || errMsg == "playoff_quarter_is_already_prepared" || errMsg == "playoff_quarter_is_already_started" {
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
	} else if stage == "semi" {
		err := ctrl.PreparePlayoffSemi(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "division_is_not_started" || errMsg == "playoff_semi_is_already_prepared" || errMsg == "playoff_semi_is_already_started" {
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
	} else if stage == "final" {
		err := ctrl.PreparePlayoffFinal(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "division_is_not_started" || errMsg == "playoff_final_is_already_prepared" || errMsg == "playoff_final_is_already_started" {
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
	}

	c.JSON(http.StatusOK, gin.H{})
}

func StartPlayoff(c *gin.Context) {

	var bodyData map[string]interface{}
	err := extractBody(c, &bodyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "invalid_body_params",
		})
		return
	}

	var stage string
	if bodyData["stage"] != nil {
		stage = bodyData["stage"].(string)
	}

	if stage == "quarter" {
		err := ctrl.StartPlayoffQuarter(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "playoff_quarter_is_not_prepared" || errMsg == "playoff_quarter_is_already_started" {
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
	} else if stage == "semi" {
		err := ctrl.StartPlayoffSemi(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "playoff_semi_is_not_prepared" || errMsg == "playoff_semi_is_already_started" {
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
	} else if stage == "final" {
		err := ctrl.StartPlayoffFinal(c.Request.Context())
		if err != nil {

			errMsg := err.Error()

			if errMsg == "playoff_final_is_not_prepared" || errMsg == "playoff_final_is_already_started" {
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
	}

	c.JSON(http.StatusOK, gin.H{})
}
