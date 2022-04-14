package data

import "gorm.io/gorm"

type Organisation struct {
	gorm.Model         `json:"-"`
	OrganisationName   string
	OrganisationNumber string `gorm:"primaryKey"`
	AddressLine1       string
	AddressLine2       string
	AddressLine3       string
	AddressLine4       string
	Town               string
	Postcode           string
}
