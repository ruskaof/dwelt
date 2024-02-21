package auth

import (
	"flag"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"time"
)

func GenerateToken() string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": jwt.TimeFunc().Add(time.Duration(*expirationSeconds) * time.Second).Unix(),
		},
	)
	tokenString, _ := token.SignedString([]byte(*key))

	return tokenString
}

func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(*key), nil
	})
	if err != nil {
		slog.Debug("error validating token: ", err)
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expirationTime.Before(time.Now()) {
		return false, nil
	}

	return token.Valid, nil
}

var key = flag.String("key", "secret", "jwt secret key")
var expirationSeconds = flag.Int("expiration", 3600, "token expiration time in seconds")
