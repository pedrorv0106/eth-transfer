package models


import (
	"github.com/jinzhu/gorm"
)

type FeeTransaction struct {
	gorm.Model

	Txid	string
	FeeTxid		string
	
	Amount	string
	State	string
	PrivKey	string
}

