package command

type RegisterOwner struct {
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
}
