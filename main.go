package main

import (
	"service-discovery/controllers"
	"service-discovery/database"
	"service-discovery/env"
	"service-discovery/models"
	"service-discovery/routes"
)

func main() {

	port := env.GetEnvironmentVariable("PORT")
	url := env.GetEnvironmentVariable("MONGO_URL")

	database.ConnectToMongoDB(models.MongoCall{
		DBURL: url,
	})

	controllers.ExecuteCronJob()

	routes.NewRoutes()
	routes.NewRoutes().Run("localhost:" + port)
	database.DisconnectMongoDB()
}
