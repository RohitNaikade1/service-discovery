package database

import (
	"service-discovery/env"

	"go.mongodb.org/mongo-driver/mongo"
)

func Database() (db string) {
	db = env.GetEnvironmentVariable("DB")
	return db
}

func UserCollectionName() (collection string) {
	collectionName := env.GetEnvironmentVariable("USER_COLLECTION")
	return collectionName
}

func CredentialCollectionName() (collection string) {
	collectionName := env.GetEnvironmentVariable("CREDENTIAL_COLLECTION")
	return collectionName
}

func RegistrationCollectionName() (collection string) {
	collectionName := env.GetEnvironmentVariable("REGISTRATION_COLLECTION")
	return collectionName
}

func UserCollection() (collection mongo.Collection) {
	collectionName := UserCollectionName()
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}

func CredentialCollection() (collection mongo.Collection) {
	collectionName := CredentialCollectionName()
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}

func RegistrationCollection() (collection mongo.Collection) {
	collectionName := RegistrationCollectionName()
	collection = *Client.Database(Database()).Collection(collectionName)
	return collection
}

//Resource collections.
func VirtualMachinesCollection() (collection string) {
	collectionName := env.GetEnvironmentVariable("VIRTUAL_MACHINES_COLLECTION")
	return collectionName
}
