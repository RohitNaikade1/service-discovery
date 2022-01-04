package main

import (
	"service-discovery/controllers"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/middlewares"
	"service-discovery/models"
	"service-discovery/server"
)

var Logger = middlewares.Logger()

func main() {

	url := env.GetEnvironmentVariable("MONGO_URL")

	database.ConnectToMongoDB(models.MongoCall{
		DBURL: url,
	})
	controllers.ExecuteCronJob()
	server.HttpsServer()

	database.DisconnectMongoDB()

}
