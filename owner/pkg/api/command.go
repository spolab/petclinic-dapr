package api

import "github.com/go-playground/validator/v10"

type RegisterOwner struct {
	Owner
}

func (cmd *RegisterOwner) Validate() error {
	v := validator.New()
	return v.Struct(cmd)
}
