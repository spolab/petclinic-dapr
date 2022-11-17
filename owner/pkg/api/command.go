package api

import "github.com/go-playground/validator/v10"

type RegisterOwner struct {
	Id         string `json:"id" validate:"required"`
	Salutation string `json:"salutation" validate:"required"`
	Surname    string `json:"name" validate:"required"`
	Name       string `json:"surname" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"email"`
}

func (cmd RegisterOwner) Validate() error {
	validate := validator.New()
	return validate.Struct(&cmd)
}
