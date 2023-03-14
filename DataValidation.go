package main

import "fmt"

var AdminID = "1234"
var AdminEmail = "JohnDoe@gmail.com"
var AdminPassword = "P@ssword"

var PubID = "5678"
var PubEmail = "Jahmur760@gmail.com"
var PubPassword = "Passw0rd10"

func ArchiveInfo()
func adminUser(id string, email string, password string) bool {
	if id == AdminID {
		return true
	} else {
		fmt.Println("ID is incorrect")
		return false
	}

	if email == AdminEmail {
		return true
	} else {
		fmt.Println("Email is correct")
		return false
	}

	if password == AdminPassword {
		return true
	} else {
		fmt.Println("Password is correct")
		return false
	}

}
func publicUser(id string, email string, password string) bool {
	if id == PubID {
		return true
	} else {
		fmt.Println("ID is incorrect")
		return false
	}

	if email == PubEmail {
		return true
	} else {
		fmt.Println("Email is incorrect")
		return false
	}

	if password == PubPassword {
		return true
	} else {
		fmt.Println("Password is incorrect")
		return false
	}

}

func history(comment string, admin_id string, edit_made string) {

}

func form()

func main() {

}
