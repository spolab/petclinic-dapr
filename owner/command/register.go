package command

type RegisterOwner struct {
	Salutation string `json:"salutation" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Address    string `json:"address"`
	PostCode   string `json:"post_code"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Email      string `json:"email" validate:"required,email"`
}
