package event

const (
	TypeClientRegisteredV1 = "ClientRegistered/v1"
)

type ClientRegistered struct {
	Id         string `json:"id"`
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}
