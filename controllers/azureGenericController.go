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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetResourceByID(subid string, id string, resourceType string) (res []byte) {

	var data map[string]interface{}

	client := armresources.NewResourcesClient(subid, GetDefaultCred(), nil)
	var Result armresources.ResourcesGetByIDResponse
	if resourceType != "databases" && resourceType != "servers" {
		resp, err := client.GetByID(context.Background(), id, "2021-04-01", nil)
		if err != nil {
			fmt.Println("failed to obtain a response: ", err.Error())
		}

		Result = resp

	} else {
		resp, err := client.GetByID(context.Background(), id, "2021-08-01-preview", nil)
		if err != nil {
			fmt.Println("failed to obtain a response: ", err.Error())
		}
		Result = resp
	}

	response, err := json.Marshal(Result)
	if err != nil {
		fmt.Println(err.Error())
	}

	sysid := PostApi(string(response))
	json.Unmarshal(response, &data)
	p := AddSysID(data, sysid)

	filter := bson.M{"id": *Result.ID}
	database.UpdateToMongo(p, "service-discovery", resourceType, filter)

	r, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}

	return r
}

func GetListOfResources(subid string, name string) (response string) {

	cnt := 0

	var s []string
	client := armresources.NewResourcesClient(subid, GetDefaultCred(), nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ResourcesListResult.Value) == 0 {
			fmt.Println("missing payload")
		}

		for _, val := range resp.ResourcesListResult.Value {

			resourceType := strings.Split(*val.Type, "/")
			resource := strings.ToLower(resourceType[len(resourceType)-1])

			if resource == name {
				fmt.Println("Name: ", *val.Name)
				//fmt.Println("ID: ", *val.ID)
				//fmt.Println(resource)
				cnt++
				//
				response := GetResourceByID(subid, *val.ID, resource)
				s = append(s, string(response))
			}

		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}

	stringByte := "[" + strings.Join(s, " ,") + "]"

	fmt.Println(cnt)
	return stringByte
}

func GetResponseForAll(c *gin.Context, name string) {
	credsid := c.Query("credsid")
	id := helpers.SubscriptionID(credsid)
	fmt.Println("sid:", id)
	c.Data(http.StatusOK, "application/json", []byte(GetListOfResources(id, name)))

}
