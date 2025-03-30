package routes

import (
	"bots/shop/handler"

	"github.com/gin-gonic/gin"
)

func LineRoutes(router *gin.Engine) {
	LineRoutes := router.Group("")
	{
		LineRoutes.GET("/lineLogin", handler.LineLoginURLHandler)
		LineRoutes.POST("/LineAcess", handler.LineAuthHandler)
		LineRoutes.POST("/sendMessage", handler.SendMessage)
	}
}
