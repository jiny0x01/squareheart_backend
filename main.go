package main

import (
	"context"
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

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s TimeZone=Asia/Seoul",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("STORYLINK_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("STORYLINK_PW"),
		os.Getenv("SSL_MODE"))
	log.Printf("db connnection str:%s\n", connStr)
	db, err := ent.Open("postgres", connStr)
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()
	ctx := context.Background()
	dbClient := client.DB{
		Client: db,
		Ctx:    ctx,
	}
	client.SetDB(&dbClient)

	// auto migration
	if err := db.Schema.Create(dbClient.Ctx); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
		return
	}
	app := fiber.New()
	app.Post("/signup", controller.SignUp)

	app.Listen(":8080")
}
