package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

// test variables
var AdminID = "1234"
var AdminEmail = "JohnDoe@gmail.com"
var AdminPassword = "P@ssword"

var PubID = "5678"
var PubEmail = "Jahmur760@gmail.com"
var PubPassword = "Passw0rd10"
var ErrorMessage = "Invalid Login"
var Message = "Access granted"

// Admin user struct
type AdminUser struct {
	gorm.Model

	AdminID        int    "gorm: 'unique_index'"
	admin_email    string "gorm:'typevarchar(100);unique_index'" // ensures each email addresses are unique
	admin_password string
}

// public user struct
type publicUser struct {
	gorm.Model

	publicUserID    int    "gorm: 'unique_index'"
	public_email    string "gorm:'typevarchar(100);unique_index'" // ensures each email addresses are unique
	public_password string
}

// history struct
type History struct {
	gorm.Model

	comment   string
	admin_id  string
	edit_made string
}

// Affiant form struct
type AffiantForm struct {
	gorm.Model

	user_id                  int
	form_id                  int
	form_status              [4]string
	archive_status           int
	affiantFullName          string
	otherNames               string
	name_change_status       string
	social_security_num      int
	social_security_date     string
	social_security_country  string
	passport_number          int
	passport_date            string
	passport_country         string
	dob                      string
	place_of_birth           string
	nationality              string
	aquired_nationality      string
	spouse_name              string
	affiant_address          string
	residencial_phone_number int
	residencial_tax_number   int
	residencial_email        string
	created_on               string
}

// public user verification
func (p publicUser) publicUserVerification() string {
	if p.public_email == PubEmail && p.public_password == PubPassword {
		PubEmail = p.public_email
		PubPassword = p.public_password
		fmt.Println("Verification Success")
		return Message
	} else {
		fmt.Println("Email address or password is incorrect")
		return ErrorMessage
	}
}

// admin user verification
func (a AdminUser) adminUserVerification() string {
	if a.admin_email == AdminEmail && a.admin_password == AdminPassword {
		AdminEmail = a.admin_email
		AdminPassword = a.admin_password
		fmt.Println("Verification Success")
		return Message
	} else {
		fmt.Println("Email address or password is incorrect")
		return ErrorMessage
	}
}

// form verification
func (A AffiantForm) AffiamtFormVerification() {
	var Option int
	A.form_status[0] = "New"
	A.form_status[1] = "Pending"
	A.form_status[2] = "Verified"
	A.form_status[3] = "Returned"

}

var db *gorm.DB
var err error

func main() {
	// Loading environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s,", host, user, dbName, password, dbPort)

	// Opening connection to database
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")

	}

	// Close connection to database when the main function finishes
	defer db.Close()

	// Make migrations to the databse if they have not already been created
	db.AutoMigrate
}
