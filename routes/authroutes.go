package routes

import (
	controller "gin/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.GET("home", controller.RecentActionsHandler())
	incomingRoutes.GET("comment", controller.CommentHandle())
	incomingRoutes.GET("cookie", controller.GetCookie())
}
