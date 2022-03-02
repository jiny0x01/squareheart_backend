package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jiny0x01/storylink_backend/app/client"
	controller "github.com/jiny0x01/storylink_backend/app/controller/user"
	"github.com/jiny0x01/storylink_backend/ent"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error db loading .env file")
		return
	}

	var db_config string
	fmt.Sprintf(db_config, "host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("IP"),
		os.Getenv("PORT"),
		os.Getenv("STORYLINK_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("STORYLINK_PW"))
	db, err := ent.Open("postgres", db_config)
	if err != nil {
		log.Fatal(err)
		return
	}
	client.SetDB(db)
	defer db.Close()

	app := fiber.New()
	app.Post("/signup", controller.SignUp)
	app.Listen(":8080")
}
