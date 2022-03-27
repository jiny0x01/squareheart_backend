package main

import (
	"github.com/jiny0x01/storylink_backend/app/client"
	controller "github.com/jiny0x01/storylink_backend/app/controller/user"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db := client.GetDB()
	defer db.Client.Close()
	defer db.Redis.Close()
	app := fiber.New()
	//	app.Use(middleware.JWTAuth)
	app.Post("/signup", controller.SignUp)

	app.Listen(":8080")
}
