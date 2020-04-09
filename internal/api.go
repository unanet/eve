package internal

import (
	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"

	"gitlab.unanet.io/devops/eve/internal/log"

	"gitlab.unanet.io/devops/eve/internal/controller"
	"gitlab.unanet.io/devops/eve/internal/middleware"
)

func StartApi() {
	r := gin.New()
	r.Use(middleware.ApiError())
	r.Use(ginlogrus.Logger(log.Logger), gin.Recovery())
	controller.PingSetup(&r.RouterGroup)
	r.Run() // listen and serve on 0.0.0.0:8080
}
