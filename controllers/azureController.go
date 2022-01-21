package controllers

import (
	"context"
	"encoding/json"
	"strconv"

	"net/http"
	"service-discovery/database"
	"service-discovery/helpers"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func StoreInDB(data []byte) (response map[string]interface{}) {
	var data_map map[string]interface{}
	sysid := PostApi(string(data))
	json.Unmarshal(data, &data_map)
	sysid_map := AddSysID(data_map, sysid)
	return sysid_map
}

func GetDefaultCred() (cred *azidentity.DefaultAzureCredential) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	helpers.PrintError(err)
	return cred
}

func GetVM(subid string) string {
	var httpResponse []string
	resourcegroup_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)
	pager := resourcegroup_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, resourcegroup := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Resource group: " + *resourcegroup.Name)
			vm_client := armcompute.NewVirtualMachinesClient(subid, GetDefaultCred(), nil)
			pager := vm_client.List(*resourcegroup.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.VirtualMachinesListResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, vm := range resp.VirtualMachinesListResult.Value {
					Logger.Info("VM: " + *vm.Name)
					instanceview := armcompute.InstanceViewTypesInstanceView
					opt := armcompute.VirtualMachinesGetOptions{
						Expand: &instanceview,
					}
					result, err := vm_client.Get(context.TODO(), *resourcegroup.Name, *vm.Name, &opt)
					helpers.PrintError(err)
					out := helpers.Encode(result)
					sysid_map := StoreInDB(out)
					filter := bson.M{"id": *result.ID}
					database.InsertOrUpdate("virtualmachines", filter, sysid_map)
					response := helpers.Encode(sysid_map)
					httpResponse = append(httpResponse, string(response))
				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetStorageAccount(subid string) string {
	var httpResponse []string
	resourcegroup_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)
	pager := resourcegroup_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, resourcegroup := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Resourcegroup: " + *resourcegroup.Name)
			storageaccount_client := armstorage.NewStorageAccountsClient(subid, GetDefaultCred(), nil)
			pager := storageaccount_client.ListByResourceGroup(*resourcegroup.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.StorageAccountsListByResourceGroupResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, storageaccount := range resp.StorageAccountsListByResourceGroupResult.Value {
					Logger.Info("Storage account: " + *storageaccount.Name)
					result, err := storageaccount_client.GetProperties(context.Background(), *resourcegroup.Name, *storageaccount.Name, nil)
					helpers.PrintError(err)
					out := helpers.Encode(result)
					sysid_map := StoreInDB(out)
					filter := bson.M{"id": *result.ID}
					database.InsertOrUpdate("storageaccounts", filter, sysid_map)
					response := helpers.Encode(sysid_map)
					httpResponse = append(httpResponse, string(response))
				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}

	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"

	return stringByte
}

func GetNetworkInterfaces(subid string) string {
	var httpResponse []string
	client := armnetwork.NewNetworkInterfacesClient(subid, GetDefaultCred(), nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.NetworkInterfacesListAllResult.Value) == 0 {
			Logger.Error("missing payload")
		}
		for _, networkinterface := range resp.NetworkInterfacesListAllResult.Value {
			Logger.Info("Network Interface: " + *networkinterface.Name)
			out := helpers.Encode(networkinterface)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *networkinterface.ID}
			database.InsertOrUpdate("networkinterfaces", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetNetworkSecurityGroups(subid string) string {
	var httpResponse []string
	client := armnetwork.NewNetworkSecurityGroupsClient(subid, GetDefaultCred(), nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.NetworkSecurityGroupsListAllResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, networksecuritygroup := range resp.NetworkSecurityGroupsListAllResult.Value {
			Logger.Info("Network Security Group: " + *networksecuritygroup.Name)
			out := helpers.Encode(networksecuritygroup)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *networksecuritygroup.ID}
			database.InsertOrUpdate("networksecuritygroups", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetDisk(subid string) string {
	var httpResponse []string
	client := armcompute.NewDisksClient(subid, GetDefaultCred(), nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.DisksListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, disk := range resp.DisksListResult.Value {
			Logger.Info("Disk: " + *disk.Name)
			output := helpers.Encode(disk)
			sysid_map := StoreInDB(output)
			filter := bson.M{"id": *disk.ID}
			database.InsertOrUpdate("disks", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetPublicIPAddresses(subid string) string {
	var httpResponse []string
	client := armnetwork.NewPublicIPAddressesClient(subid, GetDefaultCred(), nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.PublicIPAddressesListAllResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, val := range resp.PublicIPAddressesListAllResult.Value {
			Logger.Info("Public IP Address: " + *val.Name)
			out := helpers.Encode(val)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *val.ID}
			database.InsertOrUpdate("publicipaddresses", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetResourceGroups(subid string) string {
	var httpResponse []string
	client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, resourcegroup := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Resource Group: " + *resourcegroup.Name)
			out := helpers.Encode(resourcegroup)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *resourcegroup.ID}
			database.InsertOrUpdate("resourcegroups", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetVirtualNetworks(subid string) string {
	var httpResponse []string
	client := armnetwork.NewVirtualNetworksClient(subid, GetDefaultCred(), nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.VirtualNetworksListAllResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, virtualnetwork := range resp.VirtualNetworksListAllResult.Value {
			Logger.Info("Virtual Network: " + *virtualnetwork.Name)
			out := helpers.Encode(virtualnetwork)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *virtualnetwork.ID}
			database.InsertOrUpdate("virtualnetworks", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetDatabase(subid string) string {
	var httpResponse []string
	resourcegroup_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)
	pager := resourcegroup_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, resourcegroup := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Rg name" + *resourcegroup.Name)
			server_client := armsql.NewServersClient(subid, GetDefaultCred(), nil)
			server_pager := server_client.ListByResourceGroup(*resourcegroup.Name, nil)
			for server_pager.NextPage(context.Background()) {
				server_resp := server_pager.PageResponse()
				if len(server_resp.ServersListByResourceGroupResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, server := range server_resp.ServersListByResourceGroupResult.Value {
					Logger.Info("Server Name" + *server.Name)
					database_client := armsql.NewDatabasesClient(subid, GetDefaultCred(), nil)
					database_pager := database_client.ListByServer(*resourcegroup.Name, *server.Name, nil)
					for database_pager.NextPage(context.Background()) {
						database_resp := database_pager.PageResponse()
						if len(database_resp.DatabasesListByServerResult.Value) == 0 {
							Logger.Warn("missing payload")
						}
						for _, sqldatabase := range database_resp.DatabasesListByServerResult.Value {
							Logger.Info("Database: " + *sqldatabase.Name)
							out := helpers.Encode(sqldatabase)
							sysid_map := StoreInDB(out)
							filter := bson.M{"id": *sqldatabase.ID}
							database.InsertOrUpdate("databases", filter, sysid_map)
							response := helpers.Encode(sysid_map)
							httpResponse = append(httpResponse, string(response))
						}
					}
					if err := database_pager.Err(); err != nil {
						Logger.Error(err.Error())
					}
				}
			}
			if err := server_pager.Err(); err != nil {
				Logger.Error(err.Error())
			}
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetLoadBalancers(subid string) string {
	var httpResponse []string
	client := armnetwork.NewLoadBalancersClient(subid, GetDefaultCred(), nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.LoadBalancersListAllResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, loadbalancer := range resp.LoadBalancersListAllResult.Value {
			Logger.Info("Load Balancer: " + *loadbalancer.Name)
			out := helpers.Encode(loadbalancer)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *loadbalancer.ID}
			database.InsertOrUpdate("loadbalancers", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetSubnets(subid string) string {
	var httpResponse []string
	resourcegroup_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)
	virtualnetwork_client := armnetwork.NewVirtualNetworksClient(subid, GetDefaultCred(), nil)
	subnet_client := armnetwork.NewSubnetsClient(subid, GetDefaultCred(), nil)
	pager := resourcegroup_client.List(nil)
	for pager.NextPage(context.Background()) {
		resourcegroup_resp := pager.PageResponse()
		if len(resourcegroup_resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, resourcegroup := range resourcegroup_resp.ResourceGroupsListResult.Value {
			Logger.Info("Resource Group: " + *resourcegroup.Name)
			vnet_pager := virtualnetwork_client.List(*resourcegroup.Name, nil)
			for vnet_pager.NextPage(context.Background()) {
				virtualnetwork_resp := vnet_pager.PageResponse()
				if len(virtualnetwork_resp.VirtualNetworksListResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, virtualnetwork := range virtualnetwork_resp.VirtualNetworksListResult.Value {
					Logger.Info("Virtual Network: " + *virtualnetwork.Name)
					subnet_pager := subnet_client.List(*resourcegroup.Name, *virtualnetwork.Name, nil)
					for subnet_pager.NextPage(context.Background()) {
						subnet_resp := subnet_pager.PageResponse()
						if len(subnet_resp.SubnetsListResult.Value) == 0 {
							Logger.Warn("missing payload")
						}
						for _, subnet := range subnet_resp.SubnetsListResult.Value {
							Logger.Info("Subnet: " + *subnet.Name)
							subnet_response, err := subnet_client.Get(context.TODO(), *resourcegroup.Name, *virtualnetwork.Name, *subnet.Name, nil)
							helpers.PrintError(err)
							out := helpers.Encode(subnet_response)
							sysid_map := StoreInDB(out)
							filter := bson.M{"id": *subnet.ID}
							database.InsertOrUpdate("subnets", filter, sysid_map)
							response := helpers.Encode(sysid_map)
							httpResponse = append(httpResponse, string(response))
						}
					}
				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

func GetSQLServers(subid string) string {
	var httpResponse []string
	server_client := armsql.NewServersClient(subid, GetDefaultCred(), nil)
	pager := server_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ServersListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, server := range resp.ServersListResult.Value {
			Logger.Info("Server Name" + *server.Name)
			out := helpers.Encode(server)
			sysid_map := StoreInDB(out)
			filter := bson.M{"id": *server.ID}
			database.InsertOrUpdate("servers", filter, sysid_map)
			response := helpers.Encode(sysid_map)
			httpResponse = append(httpResponse, string(response))
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	return stringByte
}

//Function without authentication - used in sync api (where resources are fetched from registrations)
func GetResponce(c *gin.Context, fn func(id string) string) {
	credsid := c.Query("credsid")
	id := helpers.SubscriptionID(credsid)
	Logger.Info("sid:" + id)
	c.Data(http.StatusOK, "application/json", []byte(fn(id)))
}

//Used in single resource api
func GetResourceResponce(c *gin.Context, fn func(id string) string) {
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	credsid := c.Query("credsid")
	if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
		id := helpers.SubscriptionID(credsid)
		Logger.Info("sid:" + id)
		c.Data(http.StatusOK, "application/json", []byte(fn(id)))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

//Generic Function
func GetResourceByID(subid string, id string, resourceType string) (res []byte) {
	client := armresources.NewResourcesClient(subid, GetDefaultCred(), nil)
	var Result armresources.ResourcesGetByIDResponse
	if resourceType != "databases" && resourceType != "servers" {
		resp, err := client.GetByID(context.Background(), id, "2021-04-01", nil)
		helpers.PrintError(err)
		Result = resp
	} else {
		resp, err := client.GetByID(context.Background(), id, "2021-08-01-preview", nil)
		helpers.PrintError(err)
		Result = resp
	}
	output := helpers.Encode(Result)
	sysid_map := StoreInDB(output)
	filter := bson.M{"id": *Result.ID}
	database.InsertOrUpdate(resourceType, filter, sysid_map)
	response := helpers.Encode(sysid_map)
	return response
}

//Generic function
func GetListOfResources(subid string, name string) (response string) {
	cnt := 0
	var httpResponse []string
	client := armresources.NewResourcesClient(subid, GetDefaultCred(), nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourcesListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}
		for _, val := range resp.ResourcesListResult.Value {
			resourceType := strings.Split(*val.Type, "/")
			resource := strings.ToLower(resourceType[len(resourceType)-1])
			if resource == name {
				Logger.Info("Name: " + *val.Name)
				response := GetResourceByID(subid, *val.ID, resource)
				httpResponse = append(httpResponse, string(response))
				cnt++
			}
		}
	}
	if err := pager.Err(); err != nil {
		Logger.Error(err.Error())
	}
	stringByte := "[" + strings.Join(httpResponse, " ,") + "]"
	countStr := strconv.Itoa(cnt)
	Logger.Info("Total : " + countStr)
	return stringByte
}

//Response of generic function
func GetResponseForAll(c *gin.Context, name string) {
	credsid := c.Query("credsid")
	id := helpers.SubscriptionID(credsid)
	Logger.Info("sid:" + id)
	c.Data(http.StatusOK, "application/json", []byte(GetListOfResources(id, name)))
}
