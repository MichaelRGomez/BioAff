// Filename: BIOAFF/backend/internal/data/publicUser.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"se_group1.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

//Password stuff

// creating a password type
type password struct {
	plaintext *string
	hash      []byte
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

//PublicUser stuff

// creaing the PublicUser model
type PublicUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

// Declaring an AnonymousUser, no id, no name, no email, no password
var AnonymousUser = &PublicUser{}

// checking if a user is anonymous
func (u *PublicUser) IsAnonymous() bool {
	return u == AnonymousUser
}

// creating our publicuser model
type PublicUserModel struct {
	DB *sql.DB
}

//CRUD functions will be made later

// Insert() - creates a new publicuser on the database
func (m PublicUserModel) Insert(p_user *PublicUser) error {
	//creating the query
	query := `
		insert into public_user (name, email, pu_password, activated)
		values ($1, $2, $3, $4)
		returning id, created_at, version
	`

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//collecting the data fields into a slice for query insertion
	args := []interface{}{p_user.Name, p_user.Email, p_user.Password.hash, p_user.Activated}

	//executing the query
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&p_user.ID, &p_user.CreatedAt, &p_user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "public_user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetByEmail() - get's the publicuser based on their email
func (m PublicUserModel) GetByEmail(email string) (*PublicUser, error) {
	query := `
		select id, name, email, pu_password, activated, created_at, version
		from public_user
		where email = $1
	`

	var p_user PublicUser

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&p_user.ID,
		&p_user.Name,
		&p_user.Email,
		&p_user.Password.hash,
		&p_user.Activated,
		&p_user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &p_user, nil
}

// Update() - updates the public_user record if necessary
func (m PublicUserModel) Update(p_user *PublicUser) error {
	query := `
		update public_user
		set name = $1, email = $2, pu_password = $3, activated = $4, version = version + 1
		where id = $5 and version = $6
		returning version
	`

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//collecting the data fields into a slice for query insertion
	args := []interface{}{p_user.Name, p_user.Email, p_user.Password, p_user.Activated, p_user.ID, p_user.Version}

	//running the query
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&p_user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "public_user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

//Validation stuff

// Validating client requests
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// validating the password
func ValidatePasswordPlaintex(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// validating our user
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

//tokens will be added later
