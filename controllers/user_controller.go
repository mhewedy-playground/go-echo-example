package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	user := user{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	if user.Username != "jon" || user.Password != "shhh!" {
		return echo.ErrUnauthorized
	}

	token, err := createToken("John snow", []string{"admin"})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func createToken(name string, roles []string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["roles"] = roles
	claims["exp"] = time.Now().Add(time.Hour * 24 * 10).Unix()
	t, err := token.SignedString([]byte("secret"))
	return t, err
}

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Home")
}

func Me(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	roles := claims["roles"].([]interface{})
	return c.String(http.StatusOK,
		fmt.Sprintf("Welcome %s, you have roles: %v\n", name, roles))
}
