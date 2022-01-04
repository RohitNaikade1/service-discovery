package routes

import (
	"github.com/gin-gonic/gin"
)

type routes struct {
	Router *gin.Engine
}

func NewRoutes() routes {

	r := routes{
		Router: gin.Default(),
	}
	r.Router.Use(gin.Recovery())

	servicediscovery := r.Router.Group("/servicediscovery")

	r.AddCloudResources(servicediscovery)
	r.AddRegistrationRoutes(servicediscovery)
	r.AddCredentialsRoutes(servicediscovery)
	r.AddUserRoutes(servicediscovery)
	r.AddCronJobRoutes(servicediscovery)
	return r
}

func (r routes) Run(addr ...string) error {
	return r.Router.Run()
}
