// Filename: BIOAFF/backend/internal/data/publicUser.go
package data

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"se_group1.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// Declaring an AnonymousUser, no id, no name, no email, no password
var AnonymousUser = &PublicUser{}

// creating a password type
type password struct {
	plaintext *string
	hash      []byte
}

// creaing the PublicUser model
type PublicUser struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password password `json:"-"`
}

// checking if a user is anonymous
func (u *PublicUser) IsAnonymous() bool {
	return u == AnonymousUser
}

// Set() - stores the hash of the plaintext pasword
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// matches() method checks if the supplied password is correct
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}

	}
	return true, nil
}

// Validating client requests
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintex(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, publicUser *PublicUser) {
	v.Check(publicUser.Name != "", "name", "must be provided")
	v.Check(len(publicUser.Name) <= 500, "name", "must not be more than 500 bytes long")
	//validate the email
	ValidateEmail(v, publicUser.Email)

	//validate the password
	if publicUser.Password.plaintext != nil {
		ValidatePasswordPlaintex(v, *publicUser.Password.plaintext)
	}

	//Ensure a hash of the password was created
	if publicUser.Password.hash == nil {
		panic("missing password hash for the user")
	}
}

// creating our publicuser model
type PublicUserModel struct {
	DB *sql.DB
}

//crud functions will be made later
