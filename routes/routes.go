package routes

import (
  "net/http"

  "restapi/controller"
  "restapi/middleware" 
  "github.com/labstack/echo/v4" 

)

func Init() *echo.Echo {
  e := echo.New()

  e.GET("/", func(c echo.Context) error {
	  return c.String(http.StatusOK, "Hello, this is echo")
  })

  e.GET("/pegawai", controller.GetAllPegawai, middleware.IsAuthenticated)

  e.POST("/pegawai", controller.InsertPegawai, middleware.IsAuthenticated)

  e.PUT("/pegawai", controller.UpdatePegawai, middleware.IsAuthenticated)

  e.DELETE("/pegawai", controller.DeletePegawai, middleware.IsAuthenticated)


  e.GET("/generate-hash/:password", controller.GenerateHashPassword)

  e.POST("/login", controller.CheckLogin)


  e.GET("test-struct-validation", controller.TestStructValidation)
  e.GET("test-variable-validation", controller.TestVariabelValidation)

  return e
}