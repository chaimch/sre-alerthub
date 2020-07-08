package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

func GetRules(c *gin.Context) {
	prom := c.Query("prom")
	id := c.Query("id")

	rules := models.RulesReceiver.Get(prom, id)

	c.JSON(http.StatusOK,
		gin.H{
			"Code": 0,
			"Msg":  "",
			"Data": rules,
		})
}
