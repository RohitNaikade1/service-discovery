package controllers

import (
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/helpers"
	"service-discovery/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetCurrentLoggedInUser(username string, password string, role string) (user models.User) {
	Logger.Debug("FUNCENTRY")
	err := database.Read(env.USER_COLLECTION, bson.M{"username": username, "password": password, "role": role}).Decode(&user)
	helpers.PrintError(err)
	Logger.Debug("FUNCEXIT")
	return user
}

func GetLoggedInUser(id string) (user models.User) {
	Logger.Debug("FUNCENTRY")
	err := database.Read(env.USER_COLLECTION, bson.M{"_id": id}).Decode(&user)
	helpers.PrintError(err)
	Logger.Debug("FUNCEXIT")
	return user
}

func VerifyParentAdmin(username string, password string, role string) (result bool) {
	Logger.Debug("FUNCENTRY")
	if username == env.ADMIN_USERNAME && password == env.ADMIN_PASSWORD && role == "admin" {
		result = true
	} else {
		result = false
	}
	Logger.Debug("FUNEXIT")
	return result
}
