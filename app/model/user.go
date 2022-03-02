package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/squareheart_backend/app/client"
)

type SignUpDTO struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func CreateUser(c *fiber.Ctx, dto *SignUpDTO) error {
	db := client.GetDB()
	err := db.User.
		Create().
		SetEmail(dto.Email).
		SetNickname(dto.Nickname).
		SetPassword(dto.Password).
		Exec(c.Context())

	if err != nil {
		return err
	}
	return nil
}
