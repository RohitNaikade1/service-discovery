package helpers

import (
	"context"
	"fmt"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/models"

	"go.mongodb.org/mongo-driver/bson"
)

func SubscriptionID(credsid string) (subid string) {
	var cred models.Credentials
	if database.ValidateCollection(env.CREDENTIAL_COLLECTION) {
		Logger.Info("Collection exists")
		if database.ValidateDocument(env.CREDENTIAL_COLLECTION, bson.M{"credsid": credsid}) {
			Logger.Info("Documents exists")
			collection := database.CredentialCollection()
			err := collection.FindOne(context.TODO(), bson.M{"credsid": credsid}).Decode(&cred)
			if err != nil {
				fmt.Println(err)
			}
			Logger.Info("Subscription id: " + cred.SubscriptionID)
		} else {
			Logger.Error("Doument Not found")
		}
	} else {
		Logger.Error("Collection not found")
	}
	subid = cred.SubscriptionID
	return subid
}

func FindByCredsID(credsid string) (cred models.Credentials) {
	collection := database.CredentialCollection()
	err := collection.FindOne(context.Background(), bson.M{"credsid": credsid}).Decode(&cred)
	if err != nil {
		Logger.Error(err.Error())
	}
	return cred
}
