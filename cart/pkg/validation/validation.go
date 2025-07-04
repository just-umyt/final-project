package validation

import (
	"github.com/go-playground/validator/v10"
)

func IsValid(strct any) error {
	validate := validator.New()

	err := validate.Struct(strct)
	if err != nil {

		return err
	}

	return nil
}
