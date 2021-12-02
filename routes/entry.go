package routes

import "github.com/gin-gonic/gin"

type routes struct {
	router *gin.Engine
}

func NewRoutes() routes {
	r := routes{
		router: gin.Default(),
	}

	servicediscovery := r.router.Group("/servicediscovery")

	r.AddCloudResources(servicediscovery)
	//r.AddCronJob(servicediscovery)
	r.AddRegistrationRoutes(servicediscovery)
	r.AddCredentialsRoutes(servicediscovery)
	r.AddUserRoutes(servicediscovery)
	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}
