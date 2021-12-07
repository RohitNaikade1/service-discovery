package database

import (
	"context"
	"fmt"
	"log"
	"service-discovery/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectToMongoDB(url models.MongoCall) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(url.DBURL))
	if err != nil {
		log.Fatal(err)
	}
	Client = mongoClient
	fmt.Println("Connection established with MongoDB -", url.DBURL)
}

func DisconnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	err := Client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ended Connection to MongoDB.")
}

// Validates Collection Exists or not - returns boolean value
func ValidateCollection(db string, collection string) (result bool) {

	destination := Client.Database(db)
	filter := bson.D{}
	cursor, err := destination.ListCollectionNames(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}

	result = false
	for _, c := range cursor {
		if c == collection {
			result = true
			break
		}
	}

	return result
}

func ValidateDocument(db string, collection string, filter primitive.M) (result bool) {
	destination := Client.Database(db).Collection(collection)
	curser, err := destination.Find(context.TODO(), filter, options.Find().SetLimit(1))
	if err != nil {
		fmt.Println(err)
	}
	var results []bson.M
	er := curser.All(context.TODO(), &results)
	if er != nil {
		fmt.Println(er)
	}

	count := len(results)

	if count == 1 {
		result = true
	} else {
		result = false
	}

	return result
}

//Get All Documents
func GetAllDocuments(db string, collection string) (arr []primitive.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	destination := Client.Database(db).Collection(collection)
	cursor, err := destination.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	var results []bson.M
	e := cursor.All(ctx, &results)
	if e != nil {
		fmt.Println(e)
	}

	return results
}

func UpdateOne(db string, col string, filter primitive.M, update primitive.M) (response *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	collection := Client.Database(db).Collection(col)
	res, err := collection.UpdateOne(ctx, filter, update)
	return res, err
}

//If exists then update or else insert
func UpdateToMongo(data interface{}, db string, c string, filter primitive.M) (result *mongo.SingleResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	update := bson.M{
		"$set": data,
	}
	destination := Client.Database(db).Collection(c)
	r := destination.FindOneAndUpdate(ctx, filter, update, &opt)
	fmt.Println("Finished saving data.")
	return r, err
}

func ListCollectionNames(db string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	destination := Client.Database(db)
	filter := bson.D{{}}
	cursor, err := destination.ListCollectionNames(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}

	return cursor
}
