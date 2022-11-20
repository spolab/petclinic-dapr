package api

const MetaEventType = "X-Event-Type"

type OwnerRegistered struct {
	Id         string
	Salutation string
	Surname    string
	Name       string
	Phone      string
	Email      string
}
