package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-discovery/database"
	"service-discovery/helpers"
	"service-discovery/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegistrationExist(name string, creds string) bool {
	var result bool
	var registration models.Registration

	registration.Accounts.CredsId = creds
	registration.Name = name

	collection := database.RegistrationCollection()
	err := collection.FindOne(context.TODO(), bson.M{"credsid": registration.Accounts.CredsId, "name": registration.Name}).Decode(&registration)
	if err != nil {
		result = true
	} else {
		result = false
	}
	//fmt.Println(registration)
	return result
}

func GetRegistration(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	var registration models.Registration

	id := c.Param("id")

	collection := database.RegistrationCollection()

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if database.ValidateCollection(database.Database(), database.RegistrationCollectionName()) {
		if database.ValidateDocument(database.Database(), database.RegistrationCollectionName(), bson.M{"_id": registration.ID}) {
			if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
				c.JSON(http.StatusOK, registration)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "Registration not found"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "Data not found"})
	}
}

func GetRegistrations(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	var arr []string

	result := database.GetAllDocuments(database.Database(), database.RegistrationCollectionName())

	for _, data := range result {
		out, err := json.Marshal(data)
		if err != nil {
			Logger.Error(err.Error())
		}
		arr = append(arr, string(out))
	}

	stringByte := "[" + strings.Join(arr, " ,") + "]"

	if sysAdmin || appUser.Role == "admin" {
		c.Data(http.StatusOK, "application/json", []byte(stringByte))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

//POST
func CreateRegistration(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	//fmt.Println("Sysadmin: ", sysAdmin)
	var registration models.Registration

	registration.ID = primitive.NewObjectID().Hex()

	err := c.ShouldBind(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	registration.Created_At = time.Now().Local().String()
	registration.Updated_At = time.Now().Local().String()

	collection := database.RegistrationCollection()
	fmt.Println(registration.Name, " ", registration.Accounts.CredsId)

	if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
		if database.ValidateDocument(database.Database(), database.RegistrationCollectionName(), bson.M{"credsid": registration.Accounts.CredsId, "name": registration.Name}) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "credsid or name used already"})
		} else {
			result, err := collection.InsertOne(context.Background(), registration)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "inserted", "id": result.InsertedID})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, "unauthorized")
	}

}

//PUT
func UpdateRegistration(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	var registration models.Registration

	registration.ID = c.Param("id")

	err := c.ShouldBind(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	filter := bson.M{"_id": registration.ID}
	update := bson.M{"$set": registration}

	user := helpers.GetUserByCredsID(registration.Accounts.CredsId)

	collection := database.RegistrationCollection()

	if sysAdmin || appUser.Role == "admin" || appUser.ID == user.ID {
		response, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"Updated ID": response.UpsertedID, "Data": registration})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

func DeleteRegistration(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	id := c.Param("id")

	user := helpers.GetUserByCredsID(id)

	collection := database.RegistrationCollection()
	if sysAdmin || appUser.Role == "admin" || appUser.ID == user.ID {
		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		}
		c.JSON(http.StatusOK, gin.H{"status": "Deleted", "Deleted Count": result.DeletedCount})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
