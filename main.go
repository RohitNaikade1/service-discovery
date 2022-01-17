package main

import (
	"service-discovery/controllers"
	"service-discovery/database"
	"service-discovery/middlewares"
	"service-discovery/server"
)

var Logger = middlewares.Logger()

func main() {

	Logger.Debug("Starting application...")

	database.ConnectToMongoDB()

	controllers.ExecuteCronJob()

	server.HttpsServer()

	database.DisconnectMongoDB()

	Logger.Debug("Exit..")

}
