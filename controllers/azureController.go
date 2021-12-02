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
	"github.com/gin-gonic/gin"

	//"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"

	//"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
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

/*
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
					sysid := snow.PostApi(string(out))
					json.Unmarshal(out, &data)
					p := snow.AddSysID(data, sysid)

					filter := bson.M{"id": *result.ID}
					mongodb.UpdateToMongo(p, mongodb.DB, "storageaccounts", filter)

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
*/

func GetResponce(c *gin.Context, fn func(id string) string) {

	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")
	credsid := c.Query("credsid")

	user := helpers.GetUserByCredsID(credsid)

	if helpers.ValidateRole(role) || helpers.ValidateUser(username, password, role, user) {
		id := helpers.SubscriptionID(credsid)
		fmt.Println("sid:", id)
		c.Data(http.StatusOK, "application/json", []byte(fn(id)))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

}
