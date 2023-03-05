package controller

import (
	"strconv"
	"net/http"
	"restapi/models"
	"github.com/labstack/echo/v4"
)

func GetAllPegawai(c echo.Context) error {
	result, err := models.FetchAllPegawai()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error() } )
	}

	return c.JSON(http.StatusOK, result)
}

func InsertPegawai(c echo.Context) error {
	nama := c.FormValue("nama")
	alamat := c.FormValue("alamat")
	telepon := c.FormValue("telepon")

	result, err := models.StorePegawai(nama, alamat, telepon)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdatePegawai(c echo.Context) error {
	id := c.FormValue("id")
	nama := c.FormValue("nama")
	alamat := c.FormValue("alamat")
	telepon := c.FormValue("telepon")

	// convert tpye id from string to int
	conv_id, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result, err := models.UpdatePegawai(conv_id, nama, alamat, telepon)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		} )
	}

	return c.JSON(http.StatusOK, result)
}

func DeletePegawai(c echo.Context) error {
	id := c.FormValue("id")

	conv_id, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error() )
	}

	result, err := models.DeletePegawai(conv_id)

	return c.JSON(http.StatusOK, result)
}