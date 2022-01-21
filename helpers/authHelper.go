package helpers

import (
	"context"
	"service-discovery/database"
	"service-discovery/middlewares"
	"service-discovery/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var Logger = middlewares.Logger()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateUser(username string, password string, role string, user models.User) (result bool) {
	if username == user.UserName && password == user.Password && role == "admin" {
		result = true
	} else if username == user.UserName && password == user.Password && role == user.Role {
		result = true
	} else {
		result = false
	}
	return result
}

func ValidateRole(role string) (result bool) {
	if role == "admin" {
		result = true
	} else {
		result = false
	}
	return result
}

func GetUser(userId string) (user models.User) {
	col := database.UserCollection()
	err := col.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		Logger.Error(err.Error())
	}
	return user
}

func GetUserByCredsID(credsid string) (user models.User) {
	var cred models.Credentials
	collection := database.CredentialCollection()
	collection.FindOne(context.Background(), bson.M{"credsid": credsid}).Decode(&cred)
	userId := cred.User.ID
	user = GetUser(userId)
	return user
}

func VerifyAdmin(role string, username string, password string) (result bool) {
	var user models.User
	collection := database.UserCollection()
	collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if username == user.UserName && password == user.Password && role == "admin" && user.Role == "admin" {
		result = true
	} else {
		result = false
	}
	return result
}

func VerifyUser(role string, username string, password string) (result bool) {
	var user models.User
	collection := database.UserCollection()
	collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if username == user.UserName && password == user.Password && user.Role == role {
		result = true
	} else {
		result = false
	}
	return result
}

func GetTokenValues(c *gin.Context) (username string, password string, role string) {
	username = c.GetString("username")
	password = c.GetString("password")
	role = c.GetString("role")
	return username, password, role
}
