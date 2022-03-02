package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/storylink_backend/app/model"
)

func SignUp(c *fiber.Ctx) error {
	var dto model.SignUpDTO
	if err := c.BodyParser(dto); err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("email: %s\n", dto.Email)
	log.Printf("userId: %s\n", dto.Nickname)
	log.Printf("password: %s\n", dto.Password)
	err := model.CreateUser(c, &dto)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// TODO
	// transfer token
	log.Println("OK")
	return nil
}
