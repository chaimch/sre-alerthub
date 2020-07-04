package routers

import (
	"github.com/gin-gonic/gin"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/controllers"
)

var Router *gin.Engine

func init()  {
	Router = gin.Default()
	v1 := Router.Group("/api/v1")
	{
		v1.POST("/alerts", controllers.HandleAlerts)
	}
}
