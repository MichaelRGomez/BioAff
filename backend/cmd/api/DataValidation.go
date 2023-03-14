package main

import (
	"fmt"
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
	admin_email    string
	admin_password string
}

// public user struct
type publicUser struct {
	public_email    string
	public_password string
}

// history struct
type History struct {
	comment   string
	admin_id  string
	edit_made string
}

// Affiant form struct
type AffiantForm struct {
	user_id                  int
	form_id                  int
	form_status              string
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

}

func main() {

}
