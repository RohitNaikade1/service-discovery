package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"service-discovery/database"
	"service-discovery/helpers"

	"service-discovery/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	username := c.GetString("username")
	userPassword := c.GetString("password")
	role := c.GetString("role")

	sysAdmin := VerifyParentAdmin(username, userPassword, role)
	appUser := GetCurrentLoggedInUser(username, userPassword, role)

	if sysAdmin || appUser.Role == "admin" {
		var user models.User

		user.ID = primitive.NewObjectID().Hex()

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
		}

		password := user.Password

		hashPassword, err := helpers.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		user.Password = hashPassword
		user.Created_At = time.Now().Local().String()
		user.Updated_At = time.Now().Local().String()

		collection := database.UserCollection()

		if database.ValidateDocument(database.Database(), database.UserCollectionName(), bson.M{"email": user.Email, "username": user.UserName}) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists"})
		} else {
			result, err := collection.InsertOne(context.Background(), user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "inserted", "id": result.InsertedID})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, "Not authorized to add users")
	}

}

func Login(c *gin.Context) {
	var login models.LoginDetails
	var user models.User
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	collection := database.UserCollection()
	if database.ValidateDocument(database.Database(), database.UserCollectionName(), bson.M{"username": login.UserName}) {

		collection.FindOne(context.Background(), bson.M{"username": login.UserName}).Decode(&user)

		match := helpers.CheckPasswordHash(login.Password, user.Password)

		if match {
			expirationTime := time.Now().Add(time.Minute * 30)

			claims := &models.Claims{
				Username: user.UserName,
				Password: user.Password,
				Role:     user.Role,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				Logger.Error(err.Error())
			}
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
	}
}

func GetUsers(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	var arr []string

	if database.ValidateCollection(database.Database(), database.UserCollectionName()) {
		results := database.GetAllDocuments(database.Database(), database.UserCollectionName())

		for _, users := range results {
			response, err := json.Marshal(users)
			if err != nil {
				Logger.Error(err.Error())
			}
			arr = append(arr, string(response))
		}

		stringByte := "[" + strings.Join(arr, " ,") + "]"

		sysAdmin := VerifyParentAdmin(username, password, role)
		appUser := GetCurrentLoggedInUser(username, password, role)

		if sysAdmin || appUser.Role == "admin" {
			c.Data(http.StatusOK, "application/json", []byte(stringByte))
		} else {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
	}
}

func GetUser(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	var user models.User
	id := c.Param("id")
	collection := database.UserCollection()

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	if sysAdmin || appUser.Role == "admin" || appUser.ID == id {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to access this user details."})
	}
}

func UpdateUser(c *gin.Context) {
	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	var user models.User
	user.ID = c.Param("id")

	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
	}

	userPassword := user.Password
	hashPassword, err := helpers.HashPassword(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	user.Password = hashPassword

	filter := bson.M{"_id": user.ID}

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	Logger.Info("username: " + user.UserName + "password: " + user.Password + "role: " + user.Role)

	collection := database.UserCollection()
	if sysAdmin || appUser.Role == "admin" {
		update := bson.M{"$set": bson.M{"first_name": user.First_Name, "last_name": user.Last_Name, "password": user.Password, "email": user.Email, "role": user.Role, "updated_at": time.Now().Local().String()}}
		response, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		}
		c.JSON(http.StatusOK, gin.H{"Updated count": response.ModifiedCount, "Updated data": user})
	} else {
		if appUser.ID == user.ID {
			update := bson.M{"$set": bson.M{"first_name": user.First_Name, "last_name": user.Last_Name, "password": user.Password, "email": user.Email, "updated_at": time.Now().Local().String()}}
			response, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
			}
			c.JSON(http.StatusOK, gin.H{"Updated count": response.ModifiedCount, "Updated data": user})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to update user details."})
		}
	}
}

func DeleteUser(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	password := c.GetString("password")
	id := c.Param("id")

	sysAdmin := VerifyParentAdmin(username, password, role)
	appUser := GetCurrentLoggedInUser(username, password, role)

	if sysAdmin || appUser.Role == "admin" {
		collection := database.UserCollection()
		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		}
		c.JSON(http.StatusOK, gin.H{"Deleted count": result.DeletedCount, "Deleted Count": result.DeletedCount})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to delete user"})
	}

}
