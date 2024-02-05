package routes

import (
	"clubhub/src/apis/handlers"

	"github.com/labstack/echo/v4"
)

func InitializeHotelApis(route *echo.Group) {

	route.GET("/", handlers.GetAllHotels)
	route.POST("/", handlers.CreateHotel)
	route.PATCH("/", handlers.UpdateDetails)

}
