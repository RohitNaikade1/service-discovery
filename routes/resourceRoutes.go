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
	azres.GET(":name", func(c *gin.Context) {
		//auth.ValidateToken(c)
		name := c.Param("name")
		//c.Header("Token")
		//azure.ResourceResponse(name, c)
		switch name {

		case "virtualmachines":
			controllers.GetResponce(c, controllers.GetVM)
		/*case "resourcegroups":
			azure.GetResourceGroupsResponse(c)
		case "networkinterfaces":
			azure.GetNetworkInterfacesResponse(c)
		case "virtualnetworks":
			azure.GetVirtualNetworksResponse(c)
		case "networksecuritygroups":
			azure.GetNetworkSecurityGroupsResponse(c)
		case "storageaccounts":
			azure.GetStorageAccountResponse(c)
		case "disks":
			azure.GetDiskResponse(c)
		case "publicipaddresses":
			azure.GetPublicIPAddressesResponse(c)
		case "sqlservers":
			azure.GetSQLServersResponse(c)
		case "sqldatabases":
			azure.GetDatabasesResponse(c)
		case "loadbalancers":
			azure.GetLoadBalancersResponse(c)
		case "subnets":
			azure.GetSubnetsResponse(c)
		case "all":
			azure.GetAllResponse(c)*/
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource namespace or resource type"})
		}

	})
}
