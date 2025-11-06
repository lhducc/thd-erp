package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		msg := ""
		for _, err := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("Field '%s' failed on the '%s' tag\n", err.Field(), err.Tag())
		}
		return fmt.Errorf(msg)
	}
	return nil
}
