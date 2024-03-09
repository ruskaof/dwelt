package auth

import (
	"flag"
	"github.com/golang-jwt/jwt"
	"log/slog"
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
	slog.Debug("generating token using key: " + *key) // fixme remove
	tokenString, _ := token.SignedString([]byte(*key))

	return tokenString
}

func ValidateToken(tokenString string) (username string, valid bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(*key), nil
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

var key = flag.String("jwtkey", "secret", "jwt secret key")
var expirationSeconds = flag.Int("expiration", 3600, "token expiration time in seconds")
