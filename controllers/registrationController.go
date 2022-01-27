package controllers

import (
	"fmt"
	"net/http"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/helpers"
	"service-discovery/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegExist(name string, creds string) (result bool) {
	Logger.Debug("FUNCENTRY")
	var registration models.Registration
	err := database.Read(env.REGISTRATION_COLLECTION, bson.M{"name": name, "accounts.credsid": creds}).Decode(&registration)
	if err != nil {
		result = false
	} else {
		result = true
	}
	Logger.Debug("FUNCEXIT")
	return result
}

func GetRegistration(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	var registration models.Registration
	id := c.Param("id")
	if database.ValidateCollection(env.REGISTRATION_COLLECTION) {
		if database.ValidateDocument(env.REGISTRATION_COLLECTION, bson.M{"_id": id}) {
			if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
				err := database.Read(env.REGISTRATION_COLLECTION, bson.M{"_id": id}).Decode(&registration)
				helpers.PrintError(err)
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
	Logger.Debug("FUNCEXIT")
}

func GetRegistrations(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	var arr []string
	if sysAdmin || appUser.Role == "admin" {
		result := database.ReadAll(env.REGISTRATION_COLLECTION)
		for _, data := range result {
			response := helpers.Encode(data)
			arr = append(arr, string(response))
		}
		stringByte := "[" + strings.Join(arr, " ,") + "]"
		c.Data(http.StatusOK, "application/json", []byte(stringByte))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	Logger.Debug("FUNCEXIT")
}

//POST
func CreateRegistration(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	var registration models.Registration
	registration.ID = primitive.NewObjectID().Hex()
	err := c.ShouldBind(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	registration.Created_At = time.Now().Local().String()
	registration.Updated_At = time.Now().Local().String()
	fmt.Println(RegExist(registration.Name, registration.Accounts.CredsId))
	if RegExist(registration.Name, registration.Accounts.CredsId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credsid or name used already"})
	} else {
		if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
			result := database.Insert(env.REGISTRATION_COLLECTION, registration)
			c.JSON(http.StatusOK, gin.H{"status": "inserted", "id": result.InsertedID})
		} else {
			c.JSON(http.StatusUnauthorized, "unauthorized")
		}
	}
	Logger.Debug("FUNCEXIT")
}

//PUT
func UpdateRegistration(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	fmt.Println(sysAdmin)
	fmt.Println(appUser)
	var registration models.Registration
	id := c.Param("id")
	registration.ID = id
	err := c.ShouldBind(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	registration.Updated_At = time.Now().Local().String()
	filter := bson.M{"_id": registration.ID}
	update := bson.M{"$set": registration}
	user := helpers.GetUserByCredsID(registration.Accounts.CredsId)
	fmt.Println(appUser.ID)
	fmt.Println(user.ID)
	if sysAdmin || appUser.Role == "admin" || appUser.ID == user.ID {
		response := database.Update(env.REGISTRATION_COLLECTION, filter, update)
		c.JSON(http.StatusOK, gin.H{"Matched ": response.MatchedCount, "Updated ": response.ModifiedCount, "Upsert ": response.UpsertedCount, "Upserted ID": response.UpsertedID, "Data": registration})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	Logger.Debug("FUNCEXIT")
}

func DeleteRegistration(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	id := c.Param("id")
	user := helpers.GetUserByCredsID(id)
	if sysAdmin || appUser.Role == "admin" || appUser.ID == user.ID {
		result := database.Delete(env.REGISTRATION_COLLECTION, bson.M{"_id": id})
		c.JSON(http.StatusOK, gin.H{"status": "Deleted", "Deleted Count": result.DeletedCount})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	Logger.Debug("FUNCEXIT")
}
