package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

func GetAllProms(c *gin.Context) {
	proms := models.PromsReceiver.GetAllProms()

	c.JSON(http.StatusOK,
		gin.H{
			"Code": 0,
			"Msg":  "",
			"Data": proms,
		})
}
