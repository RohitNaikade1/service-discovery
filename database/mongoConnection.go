package database

import (
	"context"
	"service-discovery/env"
	"service-discovery/middlewares"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Logger = middlewares.Logger()

//Mongo Connection
func ConnectToMongoDB() {
	username := "sdadmin"
	password := "servicediscoverydev"
	url := "localhost:27017"
	database := Database()
	uri := "mongodb://" + username + ":" + password + "@" + url + "/" + database
	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10 sec timeout
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		Logger.Error(err.Error())
	}
	Client = mongoClient
	Logger.Info("Connection established with MongoDB -" + uri)
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
func ValidateCollection(collection string) (result bool) {
	destination := Client.Database(env.MONGODB_DATABASE)
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
func ValidateDocument(collection string, filter primitive.M) (result bool) {
	Logger.Debug("FUNCENTRY")
	destination := Client.Database(env.MONGODB_DATABASE).Collection(collection)
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
	Logger.Debug("FUNCEXIT")
	return result
}

//List of call the collections
func ListCollectionNames() []string {
	Logger.Debug("FUNCENTRY")
	destination := Client.Database(env.MONGODB_DATABASE)
	filter := bson.D{{}}
	cursor, err := destination.ListCollectionNames(context.Background(), filter)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return cursor
}

func Insert(collection string, document interface{}) *mongo.InsertOneResult {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	result, err := mongoCollection.InsertOne(context.Background(), document)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return result
}

func Update(collection string, filter interface{}, update interface{}) *mongo.UpdateResult {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	result, err := mongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return result
}

func Delete(collection string, filter primitive.M) (result *mongo.DeleteResult) {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	result, err := mongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return result
}

func Read(collection string, filter interface{}) (result *mongo.SingleResult) {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	result = mongoCollection.FindOne(context.Background(), filter)
	Logger.Debug("FUNCEXIT")
	return result
}

func ReadAll(collection string) (arr []primitive.M) {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	cursor, err := mongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		Logger.Error(err.Error())
	}
	var results []bson.M
	err = cursor.All(context.Background(), &results)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return results
}

func ReadData(collection string, filter primitive.M) (arr []primitive.M) {
	Logger.Debug("FUNCENTRY")
	mongoCollection := Client.Database(env.MONGODB_DATABASE).Collection(collection)
	cursor, err := mongoCollection.Find(context.Background(), filter)
	if err != nil {
		Logger.Error(err.Error())
	}
	var results []bson.M
	err = cursor.All(context.Background(), &results)
	if err != nil {
		Logger.Error(err.Error())
	}
	Logger.Debug("FUNCEXIT")
	return results
}

func InsertOrUpdate(collection string, filter primitive.M, update primitive.M) (response interface{}) {
	Logger.Debug("FUNCENTRY")
	if ValidateDocument(collection, filter) {
		response = Update(collection, filter, update)
	} else {
		response = Insert(collection, update)
	}
	Logger.Debug("FUNCEXIT")
	return response
}
