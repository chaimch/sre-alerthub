package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

func HandleAlerts(c *gin.Context) {
	var alerts common.Alerts
	if err := c.ShouldBindJSON(&alerts); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			})
		return
	}

	var Receiver *models.Alerts
	Receiver.AlertsHandler(&alerts)

	c.JSON(http.StatusOK,
		gin.H{
			"alerts": alerts,
		})
}
