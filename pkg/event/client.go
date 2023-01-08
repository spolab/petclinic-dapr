package event

const (
	TypeClientRegisteredV1     = "ClientRegistered/v1"
	TypePetRegisteredV1        = "PetRegistered/v1"
	TypeAppointmentRequestedV1 = "AppointmentRequested/v1"
)

type ClientRegistered struct {
	Id         string `json:"id"`
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Version    int    `json:"version"`
}

type PetRegistered struct {
	ClientID string `json:"client_id"`
	ChipID   string `json:"chip_id"`
	Name     string `json:"name"`
	Breed    string `json:"breed"`
}

type AppointmentRequested struct {
	ClientID  string `json:"client_id"`
	ChipID    string `json:"chip_id"`
	Name      string `json:"name"`
	Breed     string `json:"breed"`
	Specialty string `json:"specialty"`
}
