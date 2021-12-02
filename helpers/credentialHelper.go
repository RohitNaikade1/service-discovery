package helpers

import (
	"context"
	"fmt"
	"service-discovery/database"
	"service-discovery/models"

	"go.mongodb.org/mongo-driver/bson"
)

func SubscriptionID(credsid string) (subid string) {
	var cred models.Credentials
	if database.ValidateCollection(database.Database(), database.CredentialCollectionName()) {
		fmt.Println("Collection exists")
		if database.ValidateDocument(database.Database(), database.CredentialCollectionName(), bson.M{"credsid": credsid}) {
			fmt.Println("Documents exists")
			collection := database.CredentialCollection()
			err := collection.FindOne(context.TODO(), bson.M{"credsid": credsid}).Decode(&cred)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Subscription id: ", cred.SubscriptionID)
		} else {
			fmt.Println("Doument Not found")
		}
	} else {
		fmt.Println("Collection not found")
	}
	subid = cred.SubscriptionID
	return subid
}
