package routes

import (
	"net/http"
	"service-discovery/controllers"

	"github.com/gin-gonic/gin"
)

func (r routes) AddCloudResources(rg *gin.RouterGroup) {
	cloudresources := rg.Group("/cloudresources")
	r.AddAzureRoutes(cloudresources)
}

func (r routes) AddAzureRoutes(rg *gin.RouterGroup) { //gin.HandlerFunc
	provider := rg.Group("/azure")
	r.AddAzureService(provider)
}

func (r routes) AddAzureService(rg *gin.RouterGroup) {
	service := rg.Group("/service")
	r.AzureController(service)

}

func (r routes) AzureController(rg *gin.RouterGroup) {
	azres := rg.Group("/")
	azres.GET(":name", controllers.Authenticate, func(c *gin.Context) {
		name := c.Param("name")
		switch name {

		case "virtualmachines":
			controllers.GetResponce(c, controllers.GetVM)
		case "resourcegroups":
			controllers.GetResponce(c, controllers.GetResourceGroups)
		case "networkinterfaces":
			controllers.GetResponce(c, controllers.GetNetworkInterfaces)
		case "virtualnetworks":
			controllers.GetResponce(c, controllers.GetVirtualNetworks)
		case "networksecuritygroups":
			controllers.GetResponce(c, controllers.GetNetworkSecurityGroups)
		case "storageaccounts":
			controllers.GetResponce(c, controllers.GetStorageAccount)
		case "disks":
			controllers.GetResponce(c, controllers.GetDisk)
		case "publicipaddresses":
			controllers.GetResponce(c, controllers.GetPublicIPAddresses)
		case "sqlservers":
			controllers.GetResponce(c, controllers.GetSQLServers)
		case "sqldatabases":
			controllers.GetResponce(c, controllers.GetDatabase)
		case "loadbalancers":
			controllers.GetResponce(c, controllers.GetLoadBalancers)
		case "subnets":
			controllers.GetResponce(c, controllers.GetSubnets)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource namespace or resource type"})
		}

	})
}
