package routes

import (
	"user_management_system/controllers"
	"user_management_system/middleware"

	"github.com/gin-gonic/gin"
)

//user routes
func UserRoutes(incomingRoutes *gin.RouterGroup){
   users := incomingRoutes.Group("/users")
   users.Use(middleware.Authenticate())

   users.GET("/profile",controllers.GetUsers())
   users.GET("/:user_id",controllers.Getuser())
}