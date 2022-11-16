package api

import "github.com/go-playground/validator/v10"

type RegisterOwner struct {
	Id         string
	Salutation string
	Surname    string
	Name       string
}

func (cmd RegisterOwner) Validate() error {
	validate := validator.New()
	return validate.Struct(&cmd)
}
