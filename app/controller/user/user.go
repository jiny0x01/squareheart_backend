package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/storylink_backend/app/client"
	"github.com/jiny0x01/storylink_backend/app/model"
)

func SignUp(c *fiber.Ctx) {
	var dto model.SignUpDTO
	if err := c.BodyParser(dto); err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("email: %s\n", dto.Email)
	log.Printf("userId: %s\n", dto.Nickname)
	log.Printf("password: %s\n", dto.Password)
	db := client.GetDB()
	err := db.CreateUser(c, &dto)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// TODO
	// transfer token
	log.Println("OK")
}
