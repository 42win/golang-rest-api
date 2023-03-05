package controller

import(
	"net/http"
	// alias path
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// validasi variabel sekaligus dengan struct, 
// required,email antara tiap komponent tdk boleh pakai spasi nanti error
type Customer struct {
	Nama 	string 	`validate:"required"`
	Email 	string 	`validate:"required,email"`
	Alamat 	string 	`validate:"required"`
	Umur	int 	`validate:"gte=17,lte=35"`
}

func TestStructValidation(c echo.Context) error {
	v := validator.New()

	cust := Customer{
		Nama: "Asw",
		Email: "asw@gmail.com",
		Alamat: "aa",
		Umur: 17, 
	}

	// multi validation
	err := v.Struct(cust)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}

func TestVariabelValidation(c echo.Context) error {
	v := validator.New()

	email := "asw@gmail.com"

	err := v.Var(email, "required,email")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Email not valid",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}
