package models

import "time"

type User struct {
	ID							uint		  `gorm:"primaryKey" json:"id"`
	Name  					string    `json:"name" validate:"required"`
	Surname 				string	  `json:"surname" validate:"required"`
	Patronymic			string	  `json:"patronymic" validate:"required"`
	Address					string	  `json:"address" validate:"required"`
	PassportNumber	string		`json:"passport_number" validate:"required,passport_number_format"`
	CreatedAt     	time.Time `json:"created_at"`
	UpdatedAt     	time.Time `json:"updated_at"`
}