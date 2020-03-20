package models


import (
	"github.com/jinzhu/gorm"
)

type Blockchain struct {
	gorm.Model

	Currency     string
	Height       int
}

