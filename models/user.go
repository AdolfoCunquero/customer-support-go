package models

import "time"

type User struct {
	Username  string    `json:"username" bson:"_id"`
	FirstName string    `json:"firstName" bson:"firstName"`
	LastName  string    `json:"lastName" bson:"lastName"`
	BirthDate time.Time `json:"birthDate" bson:"birthDate"`
	RolId     int       `json:"rolId" bson:"rolId"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	Password  string    `json:"password,omitempty" bson:"password"`
	Audit     Audit     `json:"audit,omitempty" bson:"audit"`
}
