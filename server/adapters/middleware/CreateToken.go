package middleware

import (
	"main/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(account *model.Account) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = account.Email
	claims["fullname"] = account.Fullname
	claims["role_id"] = account.RoleId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("clothing_token"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
