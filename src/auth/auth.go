package auth

import (
	"dwelt/src/config"
	"flag"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateToken(userId int64) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":    jwt.TimeFunc().Add(time.Duration(*expirationSeconds) * time.Second).Unix(),
			"userId": userId,
		},
	)
	tokenString, _ := token.SignedString([]byte(config.DweltCfg.JwtKey))

	return tokenString
}

func ValidateToken(tokenString string) (userId int64, valid bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.DweltCfg.JwtKey), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expirationTime.Before(time.Now()) {
		return
	}

	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		return
	}
	userId = int64(userIdFloat)

	valid = token.Valid
	return
}

var expirationSeconds = flag.Int("expiration", 3600*24*7, "token expiration time in seconds")
