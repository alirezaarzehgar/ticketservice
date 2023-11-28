package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/alirezaarzehgar/ticketservice/config"
)

const DATE_FORMAT = "2006-01-02"

func CreateSHA256(pass string) string {
	hashByte := sha256.Sum256([]byte(pass))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}

var EXPTIME = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))

func CreateUserToken(id uint, email, user string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprint(id),
		Issuer:    email,
		Subject:   user,
		ExpiresAt: EXPTIME,
	})
	bearer, _ := token.SignedString(config.JwtSecret())
	return bearer
}
func GetToken(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	return bearer[len("Bearer "):]
}

func GetUserId(c echo.Context) uint {
	bearer := c.Request().Header.Get("Authorization")
	token, _, _ := new(jwt.Parser).ParseUnverified(bearer[len("Bearer "):], jwt.MapClaims{})
	claims := token.Claims.(jwt.MapClaims)

	_, ok := claims["jti"]
	if !ok {
		return 0
	}

	id, _ := strconv.Atoi(claims["jti"].(string))
	return uint(id)

}
