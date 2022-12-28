package vet

type RegisterVetCommand struct {
	Id      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type RegisterVetResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}
