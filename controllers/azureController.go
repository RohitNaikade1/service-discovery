package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"service-discovery/database"
	"service-discovery/helpers"
	"strings"
	"time"

	//"github.com/Azure/azure-sdk-for-go/sdk/arm"

	//"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	//"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	//"github.com/Azure/azure-sdk-for-go/sdk/sql/armsql"
	"github.com/gin-gonic/gin"

	//"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	//"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"go.mongodb.org/mongo-driver/bson"
)

var data map[string]interface{}

func GetVM(subid string) string {
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	rg_client := armresources.NewResourceGroupsClient(subid, cred, nil)
	pager := rg_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			log.Println("Rg name: ", *val1.Name)
			vm_client := armcompute.NewVirtualMachinesClient(subid, cred, nil)
			pager := vm_client.List(*val1.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.VirtualMachinesListResult.Value) == 0 {
					fmt.Println("missing payload")
				}
				for _, val2 := range resp.VirtualMachinesListResult.Value {
					fmt.Println("************VM*************")
					log.Println("VM: ", *val2.Name)

					instanceview := armcompute.InstanceViewTypesInstanceView
					opt := armcompute.VirtualMachinesGetOptions{
						Expand: &instanceview,
					}
					result, err := vm_client.Get(context.TODO(), *val1.Name, *val2.Name, &opt)
					if err != nil {
						fmt.Println(err)
					}

					out, e := json.Marshal(result)
					if e != nil {
						fmt.Println(err)
					}
					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					p := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(p, database.Database(), "virtualmachines", filter)

					r, err := json.Marshal(p)
					if err != nil {
						fmt.Println(err)
					}
					s = append(s, string(r))

					time.Sleep(time.Millisecond * 500)

				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	stringByte := "[" + strings.Join(s, " ,") + "]"
	//log.Println("Bye from VM")
	return stringByte
}

func GetStorageAccount(subid string) string {
	var s []string
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	rg_client := armresources.NewResourceGroupsClient(subid, cred, nil)
	pager := rg_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			log.Println("Rg name: ", *val1.Name)
			sa_client := armstorage.NewStorageAccountsClient(subid, cred, nil)
			pager := sa_client.ListByResourceGroup(*val1.Name, nil)
			for pager.NextPage(context.Background()) {
				resp := pager.PageResponse()
				if len(resp.StorageAccountsListByResourceGroupResult.Value) == 0 {
					fmt.Println("missing payload")
				}
				for _, val2 := range resp.StorageAccountsListByResourceGroupResult.Value {

					fmt.Println("SA: ", *val2.Name)
					result, err := sa_client.GetProperties(context.Background(), *val1.Name, *val2.Name, nil)
					if err != nil {
						fmt.Println(err)
					}
					out, e := json.Marshal(result)
					if e != nil {
						fmt.Println(err)
					}
					sysid := PostApi(string(out))
					json.Unmarshal(out, &data)
					p := AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					database.UpdateToMongo(p, database.Database(), "storageaccounts", filter)

					r, err := json.Marshal(p)
					if err != nil {
						fmt.Println(err)
					}
					s = append(s, string(r))

					time.Sleep(time.Millisecond * 500)
				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		fmt.Println(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	//fmt.Println("Bye from storage account")
	return stringByte
}

func GetNetworkInterfaces(subid string) string {
	//fmt.Println("Hii from ni")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client := armnetwork.NewNetworkInterfacesClient(subid, cred, nil)

	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.NetworkInterfacesListAllResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val := range resp.NetworkInterfacesListAllResult.Value {
			fmt.Println("*******Network Interface**********")
			log.Println("Network Interface: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "networkinterfaces", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))

			time.Sleep(time.Millisecond * 700)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	return stringByte
}

/*
func GetStorageAccount(subid string) string {
	var s []string
	fmt.Println("Hii from storage account")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient(arm.NewDefaultConnection(cred, nil), subid)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.StorageAccountListResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.StorageAccountListResult.Value {
			fmt.Println("******Storage Account*****")
			fmt.Println("Name: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}
			sysid := snow.PostApi(string(out))
			json.Unmarshal(out, &data)
			p := snow.AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			mongodb.UpdateToMongo(p, database.Database(), "storageaccount", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))
		}
	}
	if err := pager.Err(); err != nil {
		fmt.Println(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	fmt.Println("Bye from storage account")
	return stringByte
}
*/
func GetNetworkSecurityGroups(subid string) string {
	//	fmt.Println("hii from nsg")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkSecurityGroupsClient(subid, cred, nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.NetworkSecurityGroupsListAllResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.NetworkSecurityGroupsListAllResult.Value {
			fmt.Println("********NSG*********")
			log.Println("Network Security Group: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}
			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "networksecuritygroups", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))
			time.Sleep(time.Millisecond * 150)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	stringByte := "[" + strings.Join(s, " ,") + "]"
	//	fmt.Println("Bye from nsg")
	return stringByte
}

func GetDisk(subid string) string {
	//	fmt.Println("Hii from disk")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client := armcompute.NewDisksClient(subid, cred, nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.DisksListResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.DisksListResult.Value {
			log.Println("Disk: ", *val.Name)
			fmt.Println("********Disk*********")
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "disks", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))

			time.Sleep(time.Millisecond * 800)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	stringByte := "[" + strings.Join(s, " ,") + "]"
	//fmt.Println("Bye from disks")
	return stringByte
}

func GetPublicIPAddresses(subid string) string {
	//fmt.Println("Hi from public ip")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(subid, cred, nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.PublicIPAddressesListAllResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.PublicIPAddressesListAllResult.Value {
			fmt.Println("*******Public IP**********")
			log.Println("Public IP Address: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "publicipaddresses", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))

			time.Sleep(time.Millisecond * 500)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	//fmt.Println("Bye from ip")
	return stringByte
}

func GetResourceGroups(subid string) string {
	//fmt.Println("Hi from rg")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client := armresources.NewResourceGroupsClient(subid, cred, nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.ResourceGroupsListResult.Value {
			fmt.Println("*******Resource Group*********")
			log.Println("Resource Group: ", *val.Name)
			//log.Println("Location: ", *val.Location)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "resourcegroups", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))

			time.Sleep(time.Millisecond * 500)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	//fmt.Println("Bye from rg")
	return stringByte
}

func GetVirtualNetworks(subid string) string {
	// 	fmt.Println("Hi from vnet")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualNetworksClient(subid, cred, nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.VirtualNetworksListAllResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.VirtualNetworksListAllResult.Value {
			fmt.Println("********VNET********")
			log.Println("Virtual Network: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "virtualnetworks", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))

			time.Sleep(time.Millisecond * 500)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	//log.Println("Bye from vnet")
	return stringByte
}

func GetDatabase(subid string) string {
	//fmt.Println("Hi from database")
	var s []string
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	rg_client := armresources.NewResourceGroupsClient(subid, cred, nil)

	pager := rg_client.List(nil)

	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourceGroupsListResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val1 := range resp.ResourceGroupsListResult.Value {
			log.Println("Rg name", *val1.Name)
			server_client := armsql.NewServersClient(subid, cred, nil)
			pager1 := server_client.ListByResourceGroup(*val1.Name, nil)
			for pager1.NextPage(context.Background()) {
				resp1 := pager1.PageResponse()
				if len(resp1.ServersListByResourceGroupResult.Value) == 0 {
					fmt.Println("missing payload")
				}

				for _, val2 := range resp1.ServersListByResourceGroupResult.Value {
					log.Println("Server Name", *val2.Name)
					database_client := armsql.NewDatabasesClient(subid, cred, nil)
					pager2 := database_client.ListByServer(*val1.Name, *val2.Name, nil)
					for pager2.NextPage(context.Background()) {
						resp2 := pager2.PageResponse()
						if len(resp2.DatabasesListByServerResult.Value) == 0 {
							fmt.Println("missing payload")
						}
						for _, val3 := range resp2.DatabasesListByServerResult.Value {
							fmt.Println("********* Database ******")
							log.Println("Database: ", *val3.Name)
							out, err := json.Marshal(val3)
							if err != nil {
								fmt.Println(err)
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							p := AddSysID(data, sysid)

							filter := bson.M{"id": *val3.ID}
							database.UpdateToMongo(p, database.Database(), "sqldatabases", filter)

							r, err := json.Marshal(p)
							if err != nil {
								fmt.Println(err)
							}
							s = append(s, string(r))

							time.Sleep(time.Millisecond * 500)
						}
					}
					if err := pager2.Err(); err != nil {
						log.Fatal(err)
					}
				}
			}
			if err := pager1.Err(); err != nil {
				log.Fatal(err)
			}
		}

	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	stringByte := "[" + strings.Join(s, " ,") + "]"
	//	fmt.Println("Bye from database")
	return stringByte
}

func GetLoadBalancers(subid string) string {
	//fmt.Println("Hii from lb")
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client := armnetwork.NewLoadBalancersClient(subid, cred, nil)
	pager := client.ListAll(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.LoadBalancersListAllResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val := range resp.LoadBalancersListAllResult.Value {
			fmt.Println("******Load Balancer********")
			log.Println("Load Balancer: ", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}

			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)

			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "loadbalancers", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))
		}

	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	stringByte := "[" + strings.Join(s, " ,") + "]"
	return stringByte
}

func GetSubnets(subid string) string {
	var s []string

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	resourcegroup_client := armresources.NewResourceGroupsClient(subid, cred, nil)
	virtualnetwork_client := armnetwork.NewVirtualNetworksClient(subid, cred, nil)
	subnet_client := armnetwork.NewSubnetsClient(subid, cred, nil)

	pager := resourcegroup_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp1 := pager.PageResponse()
		if len(resp1.ResourceGroupsListResult.Value) == 0 {
			fmt.Println("missing payload")
		}
		for _, val1 := range resp1.ResourceGroupsListResult.Value {
			fmt.Println("*******Resource Group*********")
			log.Println("Resource Group: ", *val1.Name)
			vnet_pager := virtualnetwork_client.List(*val1.Name, nil)
			for vnet_pager.NextPage(context.Background()) {
				resp2 := vnet_pager.PageResponse()
				if len(resp2.VirtualNetworksListResult.Value) == 0 {
					fmt.Println("missing payload")
				}
				for _, val2 := range resp2.VirtualNetworksListResult.Value {
					fmt.Println("********VNET********")
					log.Println("Virtual Network: ", *val2.Name)
					subnet_pager := subnet_client.List(*val1.Name, *val2.Name, nil)
					for subnet_pager.NextPage(context.Background()) {
						subnet_resp := subnet_pager.PageResponse()
						if len(subnet_resp.SubnetsListResult.Value) == 0 {
							fmt.Println("missing payload")
						}

						for _, subnet_val := range subnet_resp.SubnetsListResult.Value {
							fmt.Println("*****Subnet*****")
							fmt.Println("Subnet: ", *subnet_val.Name)

							response, err := subnet_client.Get(context.TODO(), *val1.Name, *val2.Name, *subnet_val.Name, nil)
							if err != nil {
								fmt.Println(err)
							}

							out, err := json.Marshal(response)
							if err != nil {
								fmt.Println(err)
							}

							sysid := PostApi(string(out))
							json.Unmarshal(out, &data)
							p := AddSysID(data, sysid)

							filter := bson.M{"id": *subnet_val.ID}
							database.UpdateToMongo(p, database.Database(), "subnets", filter)

							r, err := json.Marshal(p)
							if err != nil {
								fmt.Println(err)
							}
							s = append(s, string(r))
						}
					}
				}
			}
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	return stringByte
}

func GetSQLServers(subid string) string {
	var s []string
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	server_client := armsql.NewServersClient(subid, cred, nil)
	pager := server_client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ServersListResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val := range resp.ServersListResult.Value {
			log.Println("Server Name", *val.Name)
			out, err := json.Marshal(val)
			if err != nil {
				fmt.Println(err)
			}
			sysid := PostApi(string(out))
			json.Unmarshal(out, &data)
			p := AddSysID(data, sysid)
			filter := bson.M{"id": *val.ID}
			database.UpdateToMongo(p, database.Database(), "sqlservers", filter)

			r, err := json.Marshal(p)
			if err != nil {
				fmt.Println(err)
			}
			s = append(s, string(r))
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"
	return stringByte
}

func GetResponce(c *gin.Context, fn func(id string) string) {

	credsid := c.Query("credsid")

	id := helpers.SubscriptionID(credsid)
	fmt.Println("sid:", id)
	c.Data(http.StatusOK, "application/json", []byte(fn(id)))

}

func GetResourceResponce(c *gin.Context, fn func(id string) string) {
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	credsid := c.Query("credsid")

	//user := helpers.GetUserByCredsID(credsid)

	if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
		id := helpers.SubscriptionID(credsid)
		fmt.Println("sid:", id)
		c.Data(http.StatusOK, "application/json", []byte(fn(id)))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

}
