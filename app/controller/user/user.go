package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	util "github.com/jiny0x01/storylink_backend/app/internal"
	"github.com/jiny0x01/storylink_backend/app/model"
)

func SignUp(c *fiber.Ctx) error {
	dto := new(model.SignUpDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	log.Println(dto)
	err := model.CreateUser(c, dto)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(err.Error())
	}

	// change dto.Email to uuid
	t, err := util.GenereateToken(dto.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	// regist token to redis
	util.RegistToken()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": t})
}

func SignIn(c *fiber.Ctx) error {

}
