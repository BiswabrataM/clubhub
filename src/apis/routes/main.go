package routes

import (
	"clubhub/src/utils/dtos"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitializeApis(app *echo.Echo) {
	log.Println("APIs: initializing API routes")

	v1Routes := app.Group("/apis/v1")

	v1HotelRoutes := v1Routes.Group("/hotelchain")
	InitializeHotelApis(v1HotelRoutes)

	log.Println("APIs: initializing /apis/v1/health API route")
	v1Routes.GET("/health", func(c echo.Context) error {
		data := dtos.SuccessResponse{Message: "All OK!", Data: nil}
		return c.JSON(http.StatusOK, data)
	})
	log.Println("APIs: initialized API routes")
}
