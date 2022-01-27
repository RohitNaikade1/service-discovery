package routes

import (
	"service-discovery/controllers"

	"github.com/gin-gonic/gin"
)

func (r routes) AddCredentialsRoutes(rg *gin.RouterGroup) {
	creds := rg.Group("/")
	{
		creds.GET("credentials", controllers.Authenticate, controllers.GetAllCredentials)
		creds.GET("credentials/:credsid", controllers.Authenticate, controllers.GetCredential)
		creds.POST("credentials", controllers.Authenticate, controllers.CreateCredentials)
		creds.PUT("credentials/:id", controllers.Authenticate, controllers.UpdateCredentials)
		creds.DELETE("credentials/:id", controllers.Authenticate, controllers.DeleteCredentials)
		creds.GET("credentials/data", controllers.Authenticate, controllers.GetUserCredentials)
	}

}
