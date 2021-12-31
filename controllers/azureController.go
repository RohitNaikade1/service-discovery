package controllers

import (
	"context"
	"encoding/json"
	"strconv"

	"net/http"
	"service-discovery/database"
	"service-discovery/helpers"
	logger "service-discovery/middlewares"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var data map[string]interface{}

func GetDefaultCred() (cred *azidentity.DefaultAzureCredential) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		Logger.Error(err.Error())
	}
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
					if err != nil {
						Logger.Error(err.Error())
					}

					out, e := json.Marshal(result)
					if e != nil {
						Logger.Error(err.Error())
					}

					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					sysid_map := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(sysid_map, database.Database(), "virtualmachines", filter)

					response, err := json.Marshal(sysid_map)
					if err != nil {
						logger.Logger().Error(err.Error())
					}
					httpResponse = append(httpResponse, string(response))

					//time.Sleep(time.Millisecond * 500)

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
					if err != nil {
						Logger.Error(err.Error())
					}

					out, err := json.Marshal(result)
					if err != nil {
						Logger.Error(err.Error())
					}

					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					sysid_map := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(sysid_map, database.Database(), "storageaccounts", filter)

					response, err := json.Marshal(sysid_map)
					if err != nil {
						Logger.Error(err.Error())
					}

					httpResponse = append(httpResponse, string(response))

					time.Sleep(time.Millisecond * 500)
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

			out, err := json.Marshal(networkinterface)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *networkinterface.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "networkinterfaces", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(response))

			//time.Sleep(time.Millisecond * 700)
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

			out, err := json.Marshal(networksecuritygroup)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *networksecuritygroup.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "networksecuritygroups", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(response))

			//time.Sleep(time.Millisecond * 150)
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

			out, err := json.Marshal(disk)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *disk.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "disks", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(response))

			//time.Sleep(time.Millisecond * 800)
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

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "publicipaddresses", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(r))

			time.Sleep(time.Millisecond * 500)
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

			out, err := json.Marshal(resourcegroup)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *resourcegroup.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "resourcegroups", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(response))

			//time.Sleep(time.Millisecond * 500)
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

			out, err := json.Marshal(virtualnetwork)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *virtualnetwork.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "virtualnetworks", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(response))

			//time.Sleep(time.Millisecond * 500)
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

							out, err := json.Marshal(sqldatabase)
							if err != nil {
								Logger.Error(err.Error())
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							sysid_map := AddSysID(data, sysid)

							filter := bson.M{"id": *sqldatabase.ID}
							database.UpdateToMongo(sysid_map, database.Database(), "sqldatabases", filter)

							response, err := json.Marshal(sysid_map)
							if err != nil {
								Logger.Error(err.Error())
							}

							httpResponse = append(httpResponse, string(response))

							//time.Sleep(time.Millisecond * 500)
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

			out, err := json.Marshal(loadbalancer)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *loadbalancer.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "loadbalancers", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

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
							if err != nil {
								Logger.Error(err.Error())
							}

							out, err := json.Marshal(subnet_response)
							if err != nil {
								Logger.Error(err.Error())
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							sysid_map := AddSysID(data, sysid)

							filter := bson.M{"id": *subnet.ID}
							database.UpdateToMongo(sysid_map, database.Database(), "subnets", filter)

							response, err := json.Marshal(sysid_map)
							if err != nil {
								Logger.Error(err.Error())
							}

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
			out, err := json.Marshal(server)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			sysid_map := AddSysID(data, sysid)

			filter := bson.M{"id": *server.ID}
			database.UpdateToMongo(sysid_map, database.Database(), "sqlservers", filter)

			response, err := json.Marshal(sysid_map)
			if err != nil {
				Logger.Error(err.Error())
			}

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
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")

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

	var data map[string]interface{}

	client := armresources.NewResourcesClient(subid, GetDefaultCred(), nil)
	var Result armresources.ResourcesGetByIDResponse
	if resourceType != "databases" && resourceType != "servers" {
		resp, err := client.GetByID(context.Background(), id, "2021-04-01", nil)
		if err != nil {
			Logger.Error("failed to obtain a response: " + err.Error())
		}

		Result = resp

	} else {
		resp, err := client.GetByID(context.Background(), id, "2021-08-01-preview", nil)
		if err != nil {
			Logger.Error("failed to obtain a response: " + err.Error())
		}
		Result = resp
	}

	response, err := json.Marshal(Result)
	if err != nil {
		Logger.Error(err.Error())
	}

	sysid := PostApi(string(response))
	json.Unmarshal(response, &data)
	sysid_map := AddSysID(data, sysid)

	filter := bson.M{"id": *Result.ID}
	database.UpdateToMongo(sysid_map, "service-discovery", resourceType, filter)

	r, err := json.Marshal(sysid_map)
	if err != nil {
		Logger.Error(err.Error())
	}

	return r
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
