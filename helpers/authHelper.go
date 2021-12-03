package helpers

import (
	"context"
	"fmt"
	"service-discovery/database"
	"service-discovery/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

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
		fmt.Println(err)
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
