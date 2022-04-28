package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	util "github.com/jiny0x01/storylink_backend/app/internal"
	"github.com/jiny0x01/storylink_backend/app/models"
)

func SignUp(c *fiber.Ctx) error {
	dto := new(models.SignUpDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	log.Println(dto)
	userid, err := models.CreateUser(dto)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(err.Error())
	}
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

//	SignIn client 시나리오
/*
1. access_token, refresh_token 모두 없는 경우
	- DB에 사용자가 있으면 토큰 재발급
2. refresh_token만 있는경우
	- 사용자가 Refresh endpoint에 요청하여 access_token 재발급
3. access_token, refresh_token 모두 있는 경우
	- access_token으로 로그인
*/
func SignIn(c *fiber.Ctx) error {
	dto := new(models.SignInDTO)
	if err := c.BodyParser(dto); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	user, err := models.FindUser(dto)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusNoContent)
	}

	if err := util.CompareHash(user.Password, dto.Password); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	userid := strconv.Itoa(user.ID)
	if err := util.DeleteAuth(userid); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := util.CreateToken(userid)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := util.RegistToken(userid, token); err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("%s(%s) logined. Welcome\n", user.Nickname, user.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func SignOut(c *fiber.Ctx) error {
	ad, err := util.ExtractTokenMetadata(c.Get("Authorization"))
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userId, err := util.FindAuth(ad)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = util.DeleteAuth(userId)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user_id, err := strconv.Atoi(userId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := models.DeleteUser(user_id); err != nil {
		return c.SendStatus(fiber.StatusConflict)
	}

	log.Printf("user_id:%s SignOut. Bye~\n", userId)
	return c.SendStatus(fiber.StatusOK)
}

func Refresh(c *fiber.Ctx) error {
	// access_token이 만료되면 redis에서 access_token이 없어진다
	// 사용자가 refresh_token을 들고 있다면 Refresh 요청을 통해 redis에 사용자의 refresh_token이 있는지 확인하고 있으면 access_token을 신규 발급한다.
	mapToken := fiber.Map{}
	if err := c.BodyParser(&mapToken); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	refreshToken := mapToken["refresh_token"].(string)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("test_refresh_sign"), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
		}

		//Delete the previous Refresh Token
		delErr := util.DeleteAuth(refreshUuid)
		if delErr != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		//Create new pairs of refresh and access tokens
		user_id := claims["user_id"].(string)
		ts, createErr := util.CreateToken(user_id)
		if createErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(createErr.Error())
		}
		//save the tokens metadata to redis
		saveErr := util.RegistToken(user_id, ts)
		if saveErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(saveErr.Error())
		}
		tokens := fiber.Map{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return c.Status(fiber.StatusCreated).JSON(tokens)
	} else {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"error": "refresh expired",
		})
	}
}
