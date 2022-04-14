package data

import "gorm.io/gorm"

type Employee struct {
	gorm.Model         `json:"-"`
	OrganisationNumber string
	FirstName          string
	LastName           string
}
