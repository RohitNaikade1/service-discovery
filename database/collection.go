package database

import (
	"service-discovery/env"

	"go.mongodb.org/mongo-driver/mongo"
)

func Database() (db string) {
	db = env.MONGODB_DATABASE
	return db
}

func UserCollection() (collection mongo.Collection) {
	collectionName := env.USER_COLLECTION
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}

func CredentialCollection() (collection mongo.Collection) {
	collectionName := env.CREDENTIAL_COLLECTION
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}

func RegistrationCollection() (collection mongo.Collection) {
	collectionName := env.REGISTRATION_COLLECTION
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}
