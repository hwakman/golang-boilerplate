package main

import (
	"dev-boiler/controllers"
	"dev-boiler/services"

	"net/http"

	"github.com/labstack/echo"
)

func main() {
	db := services.GetConnection()
	mongo := services.GetMongoClient()

	e := echo.New()
	e.Static("/static", "assets")
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})

	// user routes
	uController, err := controllers.NewUserController(db, mongo)
	if err != nil {
		panic("Failed to get controller")
	}

	e.GET("/user", uController.FindUsers)
	e.POST("/user", uController.CreateUser)
	e.GET("/user/:id", uController.FindUser)
	e.PUT("/user/:id", uController.UpdateUser)
	e.DELETE("/user/:id", uController.DeleteUser)

	e.Logger.Fatal(e.Start(":8000"))
}
