package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
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

func CreateUserToken(id uint, email, user, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprint(id),
		Issuer:    email,
		Subject:   user,
		ExpiresAt: EXPTIME,
		Audience:  jwt.ClaimStrings{role},
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

func GetUserRole(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	token, _, _ := new(jwt.Parser).ParseUnverified(bearer[len("Bearer "):], jwt.MapClaims{})
	claims := token.Claims.(jwt.MapClaims)

	return claims["aud"].([]any)[0].(string)
}

type eInfo struct {
	alert string
	err   error
}

func ParseBody(c echo.Context, obj any, requireds []string, mustIgnore []string) error {
	var e eInfo
	var jsonbody []byte

	body, _ := ioutil.ReadAll(c.Request().Body)
	out := make(map[string]any)

	if body == nil {
		e = eInfo{ALERT_BAD_REQUEST, fmt.Errorf("empty request is not accepted")}
		goto failure
	}
	if err := json.Unmarshal(body, &out); err != nil {
		e = eInfo{ALERT_BAD_REQUEST, fmt.Errorf("wrong json data recieved. marshalling error: %v", err)}
		goto failure
	}

	for _, i := range mustIgnore {
		for key := range out {
			if key == i {
				delete(out, i)
			}
		}
	}

	for _, r := range requireds {
		found := false
		for key := range out {
			if key == r {
				found = true
			}
		}
		if !found {
			e = eInfo{ALERT_REQUIRED_FIELDS, fmt.Errorf("%s is required on this request", r)}
			goto failure
		}
	}

	jsonbody, _ = json.Marshal(out)
	if err := json.Unmarshal(jsonbody, obj); err != nil {
		e = eInfo{ALERT_REQUIRED_FIELDS, fmt.Errorf("wrong json data recieved. marshalling error: %v", err)}
		goto failure
	}

	return nil

failure:
	slog.Debug(e.err.Error(), "data", string(body))
	c.JSON(http.StatusBadRequest, Response{Alert: e.alert})
	return e.err
}
