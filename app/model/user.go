package model

import (
	"storylink_backend/app/client"

	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/squareheart_backend/app/client"
)

type SignUpDTO struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (db *client.DB) CreateUser(c *fiber.Ctx, dto *SignUpDTO) error {
	err := db.Client.User.
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
