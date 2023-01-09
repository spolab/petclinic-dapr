package command

type RegisterClientCommand struct {
	Salutation string `json:"salutation" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
}

// Associates a pet to the owner
type RegisterPetCommand struct {
	ChipID    string `json:"" validate:"required"`
	Name      string `json:"" validate:"required"`
	Breed     string `json:"" validate:"required"`
	BirthDate string `json:"" validate:"required"`
}

// Requests an appointment for a pet. Confirmation will include Vet and Date
type RequestAppointmentCommand struct {
	ChipID    string `json:"chip_id" validate:"required"`
	Specialty string `json:"specialty" validate:"required"`
}
