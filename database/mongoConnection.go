package database

import (
	"context"
	"service-discovery/middlewares"
	"service-discovery/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Logger = middlewares.Logger()

//Mongo Connection
func ConnectToMongoDB(url models.MongoCall) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(url.DBURL))
	if err != nil {
		Logger.Error(err.Error())
	}
	Client = mongoClient
	Logger.Info("Connection established with MongoDB -" + url.DBURL)
}

//End mongo connection
func DisconnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	err := Client.Disconnect(ctx)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Info("Ended Connection to MongoDB.")
}

// Validates Collection Exists or not - returns boolean value
func ValidateCollection(db string, collection string) (result bool) {

	destination := Client.Database(db)
	filter := bson.D{}
	cursor, err := destination.ListCollectionNames(context.TODO(), filter)
	if err != nil {
		Logger.Error(err.Error())
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

//Validates document exists or not
func ValidateDocument(db string, collection string, filter primitive.M) (result bool) {
	destination := Client.Database(db).Collection(collection)
	curser, err := destination.Find(context.TODO(), filter, options.Find().SetLimit(1))
	if err != nil {
		Logger.Error(err.Error())
	}

	var results []bson.M
	err = curser.All(context.TODO(), &results)
	if err != nil {
		Logger.Error(err.Error())
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
		Logger.Error(err.Error())
	}

	var results []bson.M
	err = cursor.All(ctx, &results)
	if err != nil {
		Logger.Error(err.Error())
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
	Logger.Info("Finished saving data.")
	return r, err
}

//List of call the collections
func ListCollectionNames(db string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	destination := Client.Database(db)
	filter := bson.D{{}}
	cursor, err := destination.ListCollectionNames(ctx, filter)
	if err != nil {
		Logger.Error(err.Error())
	}

	return cursor
}
