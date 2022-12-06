package server

type OwnerState struct {
	ID         string `json:"id"`
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
}
