package routes

import (
	controller "gin/controllers"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.Use(middleware.CheckTokenValid())
	// incomingRoutes.GET("/blog/:blog_id", controller.BlogById())
	incomingRoutes.POST("/subscribe/:blog_id", controller.Subscribe())
	incomingRoutes.POST("/unsubscribe/:blog_id", controller.Unsubscribe())
	incomingRoutes.GET("/user/activity/recent-actions", controller.Subscribedblogs())
	incomingRoutes.GET("/checksub/:blog_id", controller.CheckSub())
	incomingRoutes.GET("/checklogin", controller.CheckLogin())
}
