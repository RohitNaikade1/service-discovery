package main

import (
	"service-discovery/controllers"
	"service-discovery/database"
	"service-discovery/middlewares"
	"service-discovery/server"
)

var Logger = middlewares.Logger()

func main() {

	//url := env.GetEnvironmentVariable("MONGO_URL")

	database.ConnectToMongoDB()

	controllers.ExecuteCronJob()

	server.HttpsServer()

	database.DisconnectMongoDB()

}
