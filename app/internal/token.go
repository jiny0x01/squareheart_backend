package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jiny0x01/storylink_backend/app/client"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

func CreateToken(userid string) (*TokenDetails, error) {
	var err error
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": td.AccessUuid,
		"user_id":     userid,
		"exp":         td.AtExpires,
	}
	// TODO change HS256 to RS256. because RS256 more safety
	// TODO SignedString 부분은 서명을 환경변수로 불러와서 서명해야함
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	//	td.AccessToken = access_token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	td.AccessToken, err = at.SignedString([]byte("test_access_sign"))
	if err != nil {
		return nil, err
	}
	//t, err := token.SignedString([]byte(os.GetEnv("secretSign"))

	rtClaims := jwt.MapClaims{
		"refresh_uuid": td.RefreshUuid,
		"user_id":      userid,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("test_refresh_sign"))
	if err != nil {
		return nil, err
	}

	return td, err
}

func RegistToken(userid string, td *TokenDetails) error {
	db := client.GetDB()
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC
	rt := time.Unix(td.RtExpires, 0) // converting Unix to UTC
	now := time.Now()

	errAccess := db.Redis.Set(db.Ctx, td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := db.Redis.Set(db.Ctx, td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func ExtractToken(tokenHeader string) string {
	token := strings.Split(tokenHeader, " ")
	if len(token) == 2 {
		return token[1]
	}
	return ""
}

func VerifyToken(tokenHeader string) (*jwt.Token, error) {
	tokenString := ExtractToken(tokenHeader)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// TODO getting access token sign from env
		//		return []byte(os.Getenv("ACCESS_SECRET")), nil
		return []byte("test_access_sign"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func IsValidToken(tokenHeader string) error {
	token, err := VerifyToken(tokenHeader)
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(tokenHeader string) (*AccessDetails, error) {
	token, err := VerifyToken(tokenHeader)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims) // token.Claims를 jwt.MapClaims로 type asssertion
	if !ok {
		return nil, err
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, err
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil, err
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}

func FindAuth(auth *AccessDetails) (string, error) {
	db := client.GetDB()
	userid, err := db.Redis.Get(db.Ctx, auth.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

func DeleteAuth(uuid string) (int64, error) {
	db := client.GetDB()
	deleted, err := db.Redis.Del(db.Ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
