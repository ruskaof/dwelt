package auth

import (
	"dwelt/src/config"
	"flag"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateToken(username string) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": jwt.TimeFunc().Add(time.Duration(*expirationSeconds) * time.Second).Unix(),
			"usr": username,
		},
	)
	tokenString, _ := token.SignedString([]byte(config.DweltCfg.JwtKey))

	return tokenString
}

func ValidateToken(tokenString string) (username string, valid bool, err error) {
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

	username, ok = claims["usr"].(string)
	if !ok {
		return
	}

	valid = token.Valid
	return
}

var expirationSeconds = flag.Int("expiration", 3600, "token expiration time in seconds")