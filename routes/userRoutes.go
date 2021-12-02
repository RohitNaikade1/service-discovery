package routes

import (
	"service-discovery/controllers"

	"github.com/gin-gonic/gin"
)

func (r routes) AddUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/")
	user.POST("login", controllers.Login)
	user.POST("signup", controllers.SignUp)
	user.GET("users", controllers.Authenticate, controllers.GetUsers)
	user.GET("users/:id", controllers.Authenticate, controllers.GetUser)
	user.PUT("users/:id", controllers.Authenticate, controllers.UpdateUser)
	user.DELETE("users/:id", controllers.Authenticate, controllers.DeleteUser)
}
