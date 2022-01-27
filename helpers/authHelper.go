package helpers

import (
	"service-discovery/database"
	"service-discovery/env"
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

func GetUser(userId string) (user models.User) {
	Logger.Debug("FUNCENTRY")
	err := database.Read(env.USER_COLLECTION, bson.M{"_id": userId}).Decode(&user)
	PrintError(err)
	Logger.Debug("FUNCEXIT")
	return user
}

func GetUserByCredsID(credsid string) (user models.User) {
	Logger.Debug("FUNCENTRY")
	var cred models.Credentials
	err := database.Read(env.CREDENTIAL_COLLECTION, bson.M{"credsid": credsid}).Decode(&cred)
	PrintError(err)
	userId := cred.User.ID
	user = GetUser(userId)
	Logger.Debug("FUNCEXIT")
	return user
}

func GetTokenValues(c *gin.Context) (username string, password string, role string) {
	username = c.GetString("username")
	password = c.GetString("password")
	role = c.GetString("role")
	return username, password, role
}
