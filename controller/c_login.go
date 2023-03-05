package controller

import (
	"net/http"
	"time"
	"restapi/helpers"
	"restapi/models"
	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("password")
  
	hash, _ := helpers.HashPassword(password)
  
	return c.JSON(http.StatusOK, hash)
}

func CheckLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
  
	res, err := models.CheckLogin(username, password)
	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"messages": err.Error(),
		})
	}

	if !res {
		return echo.ErrUnauthorized
	}

	// return c.String(http.StatusOK, "berhasil login")
	
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["level"] = "application"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"messages": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

}