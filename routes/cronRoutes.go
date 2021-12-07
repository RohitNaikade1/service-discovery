package routes

import (
	"service-discovery/controllers"

	"github.com/gin-gonic/gin"
)

func (r routes) AddCronJobRoutes(rg *gin.RouterGroup) {
	cj := rg.Group("/")
	cj.POST("setjob", controllers.Authenticate, controllers.SetJob)
	cj.POST("getjob", controllers.Authenticate, controllers.Task)
}
