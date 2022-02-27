package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/storylink_backend/app/controller/auth"
)

func main() {
	app := fiber.New()
	app.Post("/signup", auth.SignUp)

	app.Listen(":8080")
}
