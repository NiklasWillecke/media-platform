package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		userID: strconv.Itoa(userID),
		exp:    time.Now(),
	})
}
