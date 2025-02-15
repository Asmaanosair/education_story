package main

import (
	"integration/config"
	"integration/controllers"
	"integration/routes"
	"os"
)

func main() {
	config.ConnectDatabase()
	controllers.LoadKey()
	router := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
