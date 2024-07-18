package main

import (
	"em-test/config"
	"em-test/router"
	"em-test/validators"

	"github.com/go-playground/validator/v10"

	_ "em-test/docs"
)


func main() {
	config.InitDB()

	validate := validator.New()
	validate.RegisterValidation("passport_number_format", validators.ValidatePassportNumberFormat)

	r := router.SetupRouter(validate)

	r.Run(":8080")
}