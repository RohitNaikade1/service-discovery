package controllers

import (
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

func GetAllCredentials(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	var arr []string
	result := database.ReadAll(env.CREDENTIAL_COLLECTION)
	for _, creds := range result {
		response := helpers.Encode(creds)
		arr = append(arr, string(response))
	}
	stringByte := "[" + strings.Join(arr, " ,") + "]"
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	if sysAdmin || appUser.Role == "admin" {
		c.Data(http.StatusOK, "application/json", []byte(stringByte))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	Logger.Debug("FUNCEXIT")
}

func CreateCredentials(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	var cred models.Credentials
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	cred.ID = primitive.NewObjectID().Hex()
	if sysAdmin {
		cred.User.ID = "1"
	} else {
		cred.User.ID = appUser.ID
	}
	err := c.ShouldBind(&cred)
	helpers.PrintError(err)
	dt := time.Now().Local()
	str := cred.Provider + "-" + dt.Format("02012006150405")
	cred.CredsID = str
	cred.Created_At = dt.String()
	cred.Updated_At = time.Now().Local().String()
	if cred.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.SubscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.TenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else {
		if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
			result := database.Insert(env.CREDENTIAL_COLLECTION, cred)
			c.JSON(http.StatusOK, gin.H{"status": "inserted", "inserted id": result.InsertedID})
			Logger.Info("Inserted")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		}
	}
	Logger.Debug("FUNCEXIT")
}

func GetCredential(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	var cred models.Credentials
	id := c.Param("credsid")
	Logger.Info("Param: " + id)
	if database.ValidateCollection(env.CREDENTIAL_COLLECTION) {
		if database.ValidateDocument(env.CREDENTIAL_COLLECTION, bson.M{"_id": id}) {
			if sysAdmin || appUser.Role == "admin" || appUser.Role == "user" {
				err := database.Read(env.CREDENTIAL_COLLECTION, bson.M{"_id": id}).Decode(&cred)
				helpers.PrintError(err)
				c.JSON(http.StatusOK, cred)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "Credentials not found"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "Data not found"})
	}
	Logger.Debug("FUNCEXIT")
}

func UpdateCredentials(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	var cred models.Credentials
	id := c.Param("id")
	credentials := helpers.GetCredentialData(id)
	Logger.Info(credentials.ID)
	if database.ValidateCollection(env.CREDENTIAL_COLLECTION) {
		if database.ValidateDocument(env.CREDENTIAL_COLLECTION, bson.M{"_id": id}) {
			err := c.ShouldBind(&cred)
			helpers.PrintError(err)
			cred.User.ID = credentials.User.ID
			filter := bson.M{"_id": id}
			update := bson.M{"$set": bson.M{"provider": cred.Provider, "subscriptionid": cred.SubscriptionID, "tenantid": cred.TenantID, "username": cred.UserName, "updated_at": time.Now().Local().String()}}
			if sysAdmin || appUser.Role == "admin" || appUser.ID == credentials.User.ID {
				response := database.Update(env.CREDENTIAL_COLLECTION, filter, update)
				c.JSON(http.StatusOK, gin.H{"Updated count": response.ModifiedCount, "Updated data": cred})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "Credentials not found"})
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "Data not found"})
	}
}

func DeleteCredentials(c *gin.Context) {
	Logger.Debug("FUNCENTRY")
	username, password, role := helpers.GetTokenValues(c)
	id := c.Param("id")
	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)
	credentials := helpers.GetCredentialData(id)
	user := helpers.GetUser(credentials.User.ID)
	if sysAdmin || appUser.Role == "admin" || appUser.Role == user.Role {
		Logger.Info("ID: " + credentials.ID + " CredsID: " + credentials.CredsID)
		result := database.Delete(env.CREDENTIAL_COLLECTION, bson.M{"_id": id})
		c.JSON(http.StatusOK, gin.H{"Deleted Count": result.DeletedCount, "Data": credentials})
	} else {
		Logger.Info("Not authorized")
		c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	Logger.Debug("FUNCEXIT")
}
