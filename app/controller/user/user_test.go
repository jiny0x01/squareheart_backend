package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	controller "github.com/jiny0x01/storylink_backend/app/controller/user"
	"github.com/jiny0x01/storylink_backend/app/model"
)

func TestSignUp(t *testing.T) {
	tests := []struct {
		description string
		param       model.SignUpDTO
		statusCode  int
	}{
		{
			description: "case 1: correct",
			param: model.SignUpDTO{
				Email:    "signup_test1@test.com",
				Nickname: "signup_test1",
				Password: "test1234!!",
			},
			statusCode: fiber.StatusOK,
		},
		{
			description: "case 2: already exist email",
			param: model.SignUpDTO{
				Email:    "signup_test1@test.com",
				Nickname: "signup_test1",
				Password: "test1234!!",
			},
			statusCode: fiber.StatusConflict,
		},
		{
			description: "case 3: invalid email type",
			param: model.SignUpDTO{
				Email:    "signup_test3",
				Nickname: "signup_test3",
				Password: "test1234!!",
			},
			statusCode: fiber.StatusUnprocessableEntity,
		},
		{
			description: "case 4: Too short password",
			param: model.SignUpDTO{
				Email:    "signup_test4@test.com",
				Nickname: "signup_test4",
				Password: "1",
			},
			statusCode: fiber.StatusUnprocessableEntity,
		},
	}

	app := fiber.New()
	app.Post("/signup", controller.SignUp)

	for _, test := range tests {
		req := httptest.NewRequest("POST", "http://localhost:1513", nil)
		req.Header.Set("X-Custom-Header", "hi")
		res, _ := app.Test(req, 1)
		if res.StatusCode != test.statusCode {
			t.Errorf("%s failed", test.description)
		}
	}
}
