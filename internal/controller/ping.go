package controller

import "github.com/gin-gonic/gin"

func PingSetup(g *gin.RouterGroup) {
	g.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
