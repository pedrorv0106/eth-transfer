package models


import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model

	Name     string
	Address  string
	PrivKey  string
}

