package api

const MetaEventType = "petclinic.spolab.github.com/eventType"

type OwnerRegistered struct {
	Id         string
	Salutation string
	Surname    string
	Name       string
	Phone      string
	Email      string
}
