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
	userid, err := model.CreateUser(c, dto)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(err.Error())
	}
	log.Printf("userid:%s\n", userid)
	// change dto.Email to uuid
	token, err := util.CreateToken(dto.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	// regist token to redis
	err = util.RegistToken(userid, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func SignIn(c *fiber.Ctx) error {
	//token := c.Get("Authorization")
	return nil
}
