package app

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
)

func InitDB(DB_URL string) *gorm.DB{

	log.Printf("DB_URL: %v\n", DB_URL)
	d, err := gorm.Open("mysql", DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	db = d

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	return db
}

const (
	userDBSchema = "./database/schema.sql"

	updateQuery           = "update"
	updateResetTokenQuery = "update-reset-token"
	updatePasswordQuery   = "update-password"
	updateOrderQuery      = "update-order"
	updateAccountQuery    = "update-account"

	insertQuery               = "insert"
	insertAuthenticationQuery = "insert-authentication"
	insertMemberQuery         = "insert-member"
	selectMemberIdQuery       = "select-member-id"
	insertId_DocumentQuery    = "insert-Id_Document"
	insertAccountQuery        = "insert-Account"
	insertOrderQuery          = "insert-order"

	selectLoginQuery      = "select-login"
	selectEmailQuery      = "select-email"
	selectResetTokenQuery = "select-reset-token"
	selectMemberQuery     = "select-member"
	selectAccountQuery    = "select-account"

	LOCKING_BUFFER_FACTOR = "1.1"
)
