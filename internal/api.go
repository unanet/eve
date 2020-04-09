package internal

import (
	"github.com/gin-gonic/gin"
	"gitlab.unanet.io/devops/eve/internal/controllers/ping"
	"gitlab.unanet.io/devops/eve/internal/middleware"
)

func StartApi() {
	r := gin.Default()
	r.Use(middleware.ApiError())
	ping.Setup(&r.RouterGroup)
	r.Run() // listen and serve on 0.0.0.0:8080
}
