package routes

import (
	"service-discovery/controllers"

	"github.com/gin-gonic/gin"
)

func (r routes) AddRegistrationRoutes(rg *gin.RouterGroup) {
	reg := rg.Group("/")
	reg.GET("registration", controllers.Authenticate, controllers.GetRegistrations)
	reg.GET("registration/:id", controllers.Authenticate, controllers.GetRegistration)
	reg.POST("registration", controllers.Authenticate, controllers.CreateRegistration)
	reg.PUT("registration/:id", controllers.Authenticate, controllers.UpdateRegistration)
	reg.DELETE("registration/:id", controllers.Authenticate, controllers.DeleteRegistration)
}
