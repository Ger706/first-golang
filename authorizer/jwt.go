package authorizer

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtFormat struct {
	Username string
	UserId   uint
}

func CreateToken(data *JwtFormat) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data.Username,
		"id":   data.UserId,
		"exp":  time.Now().Add(time.Hour * 48).Unix(),
	})
	tokenString, err := token.SignedString([]byte("golang-testing-jwt"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
