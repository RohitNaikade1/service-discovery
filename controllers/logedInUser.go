package controllers

import (
	"context"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetCurrentLoggedInUser(username string, password string, role string) (user models.User) {
	collection := database.UserCollection()
	collection.FindOne(context.Background(), bson.M{"username": username, "password": password, "role": role}).Decode(&user)
	return user
}

func VerifyParentAdmin(username string, password string, role string) (result bool) {
	if username == env.ADMIN_USERNAME && password == env.ADMIN_PASSWORD && role == "admin" {
		result = true
	} else {
		result = false
	}
	return result
}
