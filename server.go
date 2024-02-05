package main

import (
	"clubhub/configs"

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
	routes.InitializeApis(echoApp)
	databases.InitializePgDb()

	echoApp.Start(configs.AppPort)

	log.Println("initialized app listening to 8080")
}
