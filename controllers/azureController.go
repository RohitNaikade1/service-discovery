package controllers

import (
	"context"
	"encoding/json"
	"log"
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

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Rg name: " + *val1.Name)

			vm_client := armcompute.NewVirtualMachinesClient(subid, GetDefaultCred(), nil)

			pager := vm_client.List(*val1.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.VirtualMachinesListResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, val2 := range resp.VirtualMachinesListResult.Value {
					Logger.Info("VM: " + *val2.Name)

					instanceview := armcompute.InstanceViewTypesInstanceView
					opt := armcompute.VirtualMachinesGetOptions{
						Expand: &instanceview,
					}

					result, err := vm_client.Get(context.TODO(), *val1.Name, *val2.Name, &opt)
					if err != nil {
						Logger.Error(err.Error())
					}

					out, e := json.Marshal(result)
					if e != nil {
						Logger.Error(err.Error())
					}

					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					p := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(p, database.Database(), "virtualmachines", filter)

					r, err := json.Marshal(p)
					if err != nil {
						logger.Logger().Error(err.Error())
					}
					httpResponse = append(httpResponse, string(r))

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

func GetStorageAccount(subid string) string {
	var httpResponse []string

	rg_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)

	pager := rg_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Rg name: " + *val1.Name)
			sa_client := armstorage.NewStorageAccountsClient(subid, GetDefaultCred(), nil)
			pager := sa_client.ListByResourceGroup(*val1.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.StorageAccountsListByResourceGroupResult.Value) == 0 {
					Logger.Warn("missing payload")
				}

				for _, val2 := range resp.StorageAccountsListByResourceGroupResult.Value {
					Logger.Info("SA: " + *val2.Name)

					result, err := sa_client.GetProperties(context.Background(), *val1.Name, *val2.Name, nil)
					if err != nil {
						Logger.Error(err.Error())
					}

					out, err := json.Marshal(result)
					if err != nil {
						Logger.Error(err.Error())
					}

					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					p := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(p, database.Database(), "storageaccounts", filter)

					r, err := json.Marshal(p)
					if err != nil {
						Logger.Error(err.Error())
					}

					httpResponse = append(httpResponse, string(r))

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

		for _, val := range resp.NetworkInterfacesListAllResult.Value {
			Logger.Info("Network Interface: " + *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "networkinterfaces", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(r))

			time.Sleep(time.Millisecond * 700)
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

		for _, val := range resp.NetworkSecurityGroupsListAllResult.Value {
			Logger.Info("Network Security Group: " + *val.Name)

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "networksecuritygroups", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(r))
			time.Sleep(time.Millisecond * 150)
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

		for _, val := range resp.DisksListResult.Value {
			Logger.Info("Disk: " + *val.Name)

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "disks", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(r))

			time.Sleep(time.Millisecond * 800)
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

		for _, val := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Resource Group: " + *val.Name)

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "resourcegroups", filter)

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

func GetVirtualNetworks(subid string) string {

	var httpResponse []string

	client := armnetwork.NewVirtualNetworksClient(subid, GetDefaultCred(), nil)

	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.VirtualNetworksListAllResult.Value) == 0 {
			Logger.Warn("missing payload")
		}

		for _, val := range resp.VirtualNetworksListAllResult.Value {
			Logger.Info("Virtual Network: " + *val.Name)

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "virtualnetworks", filter)

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

func GetDatabase(subid string) string {
	var httpResponse []string

	rg_client := armresources.NewResourceGroupsClient(subid, GetDefaultCred(), nil)

	pager := rg_client.List(nil)

	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			Logger.Info("Rg name" + *val1.Name)
			server_client := armsql.NewServersClient(subid, GetDefaultCred(), nil)
			pager1 := server_client.ListByResourceGroup(*val1.Name, nil)
			for pager1.NextPage(context.Background()) {
				resp1 := pager1.PageResponse()
				if len(resp1.ServersListByResourceGroupResult.Value) == 0 {
					Logger.Warn("missing payload")
				}

				for _, val2 := range resp1.ServersListByResourceGroupResult.Value {
					log.Println("Server Name", *val2.Name)
					database_client := armsql.NewDatabasesClient(subid, GetDefaultCred(), nil)
					pager2 := database_client.ListByServer(*val1.Name, *val2.Name, nil)
					for pager2.NextPage(context.Background()) {
						resp2 := pager2.PageResponse()
						if len(resp2.DatabasesListByServerResult.Value) == 0 {
							Logger.Warn("missing payload")
						}

						for _, val3 := range resp2.DatabasesListByServerResult.Value {
							Logger.Info("Database: " + *val3.Name)

							out, err := json.Marshal(val3)
							if err != nil {
								Logger.Error(err.Error())
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							p := AddSysID(data, sysid)

							filter := bson.M{"id": *val3.ID}
							database.UpdateToMongo(p, database.Database(), "sqldatabases", filter)

							r, err := json.Marshal(p)
							if err != nil {
								Logger.Error(err.Error())
							}

							httpResponse = append(httpResponse, string(r))

							time.Sleep(time.Millisecond * 500)
						}
					}

					if err := pager2.Err(); err != nil {
						Logger.Error(err.Error())
					}
				}
			}

			if err := pager1.Err(); err != nil {
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

		for _, val := range resp.LoadBalancersListAllResult.Value {
			Logger.Info("Load Balancer: " + *val.Name)

			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "loadbalancers", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}

			httpResponse = append(httpResponse, string(r))
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
		resp1 := pager.PageResponse()
		if len(resp1.ResourceGroupsListResult.Value) == 0 {
			Logger.Warn("missing payload")
		}

		for _, val1 := range resp1.ResourceGroupsListResult.Value {

			Logger.Info("Resource Group: " + *val1.Name)
			vnet_pager := virtualnetwork_client.List(*val1.Name, nil)
			for vnet_pager.NextPage(context.Background()) {
				resp2 := vnet_pager.PageResponse()
				if len(resp2.VirtualNetworksListResult.Value) == 0 {
					Logger.Warn("missing payload")
				}
				for _, val2 := range resp2.VirtualNetworksListResult.Value {
					Logger.Info("Virtual Network: " + *val2.Name)
					subnet_pager := subnet_client.List(*val1.Name, *val2.Name, nil)
					for subnet_pager.NextPage(context.Background()) {
						subnet_resp := subnet_pager.PageResponse()
						if len(subnet_resp.SubnetsListResult.Value) == 0 {
							Logger.Warn("missing payload")
						}

						for _, subnet_val := range subnet_resp.SubnetsListResult.Value {
							Logger.Info("Subnet: " + *subnet_val.Name)

							subnet_response, err := subnet_client.Get(context.TODO(), *val1.Name, *val2.Name, *subnet_val.Name, nil)
							if err != nil {
								Logger.Error(err.Error())
							}

							out, err := json.Marshal(subnet_response)
							if err != nil {
								Logger.Error(err.Error())
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							p := AddSysID(data, sysid)

							filter := bson.M{"id": *subnet_val.ID}
							database.UpdateToMongo(p, database.Database(), "subnets", filter)

							r, err := json.Marshal(p)
							if err != nil {
								Logger.Error(err.Error())
							}

							httpResponse = append(httpResponse, string(r))
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

		for _, val := range resp.ServersListResult.Value {
			Logger.Info("Server Name" + *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				Logger.Error(err.Error())
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "sqlservers", filter)

			r, err := json.Marshal(p)
			if err != nil {
				Logger.Error(err.Error())
			}
			httpResponse = append(httpResponse, string(r))
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
	p := AddSysID(data, sysid)

	filter := bson.M{"id": *Result.ID}
	database.UpdateToMongo(p, "service-discovery", resourceType, filter)

	r, err := json.Marshal(p)
	if err != nil {
		Logger.Error(err.Error())
	}

	return r
}

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
