package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nirav114/url-shortner-backend.git/config"
)

func CreateJWT(secret []byte, userID int64) (string, error) {
	expiration := time.Second * time.Duration(config.EnvConfig.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":    strconv.FormatInt(userID, 10),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
