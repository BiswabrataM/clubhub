package main

import (
	"clubhub/src/apis/routes"
	"clubhub/src/databases"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Println("initializing echo server app")

	echoApp := echo.New()

	log.Println("Logger: initializing echo logger")
	echoApp.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: os.Stdout,
	}))
	log.Println("Logger: initialized echo logger")
	routes.InitializeApis(echoApp)
	databases.InitializePgDb()

	echoApp.Start(":8080")

	log.Println("initialized app listening to 8080")
}
