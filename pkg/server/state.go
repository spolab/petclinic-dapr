package server

type OwnerState struct {
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Address    string `json:"address"`
	PostCode   string `json:"post_code"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}
