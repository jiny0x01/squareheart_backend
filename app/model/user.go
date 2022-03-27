package model

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jiny0x01/storylink_backend/app/client"
	util "github.com/jiny0x01/storylink_backend/app/internal"
)

type SignUpDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Nickname string `json:"nickname" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func CreateUser(c *fiber.Ctx, dto *SignUpDTO) (string, error) {
	pw, err := util.HashPassword(dto.Password)
	if err != nil {
		return "", err
	}

	db := client.GetDB()
	user, err := db.Client.User.
		Create().
		SetEmail(dto.Email).
		SetNickname(dto.Nickname).
		SetPassword(pw).
		Save(db.Ctx)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(user.ID), nil
}
