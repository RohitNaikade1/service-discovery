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

func GetAllCredentials(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")
	var arr []string
	result := database.GetAllDocuments(database.Database(), database.CredentialCollectionName())
	for _, creds := range result {
		out, err := json.Marshal(creds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		arr = append(arr, string(out))
	}
	stringByte := "[" + strings.Join(arr, " ,") + "]"
	if helpers.VerifyAdmin(role, username, password) {
		c.Data(http.StatusOK, "application/json", []byte(stringByte))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

func CreateCredentials(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")
	var cred models.Credentials

	cred.ID = primitive.NewObjectID().Hex()

	err := c.ShouldBind(&cred)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
	}

	dt := time.Now().Local()
	str := cred.Provider + "-" + dt.Format("02012006150405")
	cred.CredsID = str
	cred.Created_At = dt.String()
	cred.Updated_At = time.Now().Local().String()

	user := helpers.GetUser(cred.User.ID)

	if cred.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.SubscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.TenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else if cred.User.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "error": "Enter all the required details"})
	} else {
		collection := database.CredentialCollection()
		if helpers.VerifyAdmin(role, username, password) || helpers.ValidateUser(username, password, role, user) {
			result, err := collection.InsertOne(context.Background(), cred)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "inserted", "id": result.InsertedID})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
		}
	}
}

func GetCredential(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	var cred models.Credentials

	cred.CredsID = c.Param("credsid")

	collection := database.CredentialCollection()
	if database.ValidateCollection(database.Database(), database.CredentialCollectionName()) {
		if database.ValidateDocument(database.Database(), database.CredentialCollectionName(), bson.M{"credsid": cred.CredsID}) {
			err := collection.FindOne(context.TODO(), bson.M{"credsid": cred.CredsID}).Decode(&cred)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			user := helpers.GetUser(cred.User.ID)
			fmt.Println(user.First_Name)

			if helpers.ValidateRole(role) || helpers.ValidateUser(username, password, role, user) {
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
}

func UpdateCredentials(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	var cred models.Credentials

	cred.CredsID = c.Param("credsid")

	if database.ValidateCollection(database.Database(), database.CredentialCollectionName()) {
		if database.ValidateDocument(database.Database(), database.CredentialCollectionName(), bson.M{"credsid": cred.CredsID}) {
			err := c.ShouldBind(&cred)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(cred)
			user := helpers.GetUser(cred.User.ID)
			filter := bson.M{"credsid": cred.CredsID}
			update := bson.M{"$set": bson.M{"provider": cred.Provider, "subscriptionid": cred.SubscriptionID, "tenantid": cred.TenantID, "username": cred.UserName, "updated_at": time.Now().Local().String()}}

			collection := database.CredentialCollection()

			if helpers.ValidateRole(role) || helpers.ValidateUser(username, password, role, user) {
				response, err := collection.UpdateOne(context.Background(), filter, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					c.Abort()
				}
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
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")

	id := c.Param("id")

	user := helpers.GetUser(id)

	collection := database.CredentialCollection()
	if helpers.ValidateRole(role) || helpers.ValidateUser(username, password, role, user) {
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
