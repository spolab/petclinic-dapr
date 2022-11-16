package api

type Name struct {
}

// swagger:
type Owner struct {
	Id         string `json:"id" validation:"required"`
	Salutation string `json:"salutation" validation:"required"`
	Surname    string `json:"surname" validation:"required"`
	Name       string `json:"name" validation:"required"`
}
