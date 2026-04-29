package routes

import (
	"user_management_system/controllers"
	"user_management_system/middleware"
	"user_management_system/notifier"

	"github.com/gin-gonic/gin"
)

//authentication routes like /api/v1/auth/signup
func AuthRoutes(incomingRoutes *gin.RouterGroup,n *notifier.Notifier){

	auth:=incomingRoutes.Group("/auth")

    auth.POST("/signup",controllers.Signup(n))
    auth.POST("/login",controllers.Login())
	auth.DELETE("/delete/:id",middleware.Authenticate(), controllers.DeleteUser())
}
