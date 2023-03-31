// Filename: BIOAFF/backend/internal/data/form.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// creating the form type
type Form struct {
	PublicUser_ID          int64     `json:"puid"`
	ID                     int64     `json:"id"`
	Status                 string    `json:"status"`
	Archive                bool      `json:"archive"`
	Fullname               string    `json:"fullname"`
	Othernames             string    `json:"othername"`
	Has_Changed_Name       bool      `json:"changed_name"`
	SocialSecurity_Number  int       `json:"ssnumber"`
	SocialSecurity_Date    string    `json:"ssdate"`
	SocialSecurity_Country string    `json:"sscountry"`
	Passport_Number        string    `json:"passport_number"`
	Passport_Date          string    `json:"passport_date"`
	Passport_Country       string    `json:"passport_country"`
	DOB                    string    `json:"dob"`
	Place_of_Birth         string    `json:"place_of_birth"`
	Nationality            string    `json:"nationality"`
	Acquired_Nationality   string    `json:"acquired_nationality"`
	Spouse_Name            string    `json:"spouse_name"`
	Address                string    `json:"address"`
	Phone_Number           string    `json:"phone"`
	Fax_Number             string    `json:"fax"`
	Residential_Email      string    `json:"residential_rmail"`
	CreatedAt              time.Time `json:"created_at"`
	Version                int32     `json:"version"`
}

//Validations will go here later CRUD functions were priority

// creatomg the FormModel for the sql.DB
type FormModel struct {
	DB *sql.DB
}

//CRUD functions for the FromModel

// Insert() - creates a new record of a form on the database;
// assuming the info from the model is clean
func (m FormModel) Insert(form *Form) error {
	query := `
		INSERT INTO form (user_id, form_status, archive_status, full_name, other_names, name_change_status,
			social_security_num, social_security_date, social_security_country, passport_number, passport_date, passport_country,
			dob, place_of_birth, nationality, acquired_nationality, spouse_name, address, residential_phone_number, residential_fax_number,
			residential_email)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING form_id, created_on, version
	`

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//collecting the data fields into a slice for query insertion
	args := []interface{}{
		form.PublicUser_ID, form.Status, form.Archive, form.Fullname, form.Othernames, form.Has_Changed_Name, form.SocialSecurity_Number,
		form.SocialSecurity_Date, form.SocialSecurity_Country, form.Passport_Number, form.Passport_Date, form.Passport_Country,
		form.DOB, form.Place_of_Birth, form.Nationality, form.Acquired_Nationality, form.Spouse_Name, form.Address, form.Phone_Number,
		form.Fax_Number, form.Residential_Email,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&form.ID, &form.CreatedAt, &form.Version)
}

// Get() - allows us to retrieve infromation from a form based on the form_id supplied
func (m FormModel) Get(id int64) (*Form, error) {
	//Ensuring that the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//constructing our query based on the id
	query := `
		SELECT user_id, form_id, form_status, archive_status, full_name, other_names, name_change_status,
		social_security_num, social_security_date, social_security_country, passport_number, passport_date, passport_country,
		dob, place_of_birth, nationality, acquired_nationality, spouse_name, address, residential_phone_number, residential_fax_number,
		residential_email, created_on, version
		FROM form
		WHERE form_id = $1
	`

	//Creating a form variable to hold the data
	var form Form

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Executing the query, and limiting the query execution to 3 seconds
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&form.PublicUser_ID,
		&form.ID,
		&form.Status,
		&form.Archive,
		&form.Fullname,
		&form.Othernames,
		&form.Has_Changed_Name,
		&form.SocialSecurity_Number,
		&form.SocialSecurity_Date,
		&form.SocialSecurity_Country,
		&form.Passport_Number,
		&form.Passport_Date,
		&form.Passport_Country,
		&form.DOB,
		&form.Place_of_Birth,
		&form.Nationality,
		&form.Acquired_Nationality,
		&form.Spouse_Name,
		&form.Address,
		&form.Phone_Number,
		&form.Fax_Number,
		&form.Residential_Email,
		&form.CreatedAt,
		&form.Version,
	)

	//in case of any errors
	if err != nil {
		//checking the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//success
	return &form, nil

}

// Update() - to edit/alter a specific form
// optimistic locking (version number)
func (m FormModel) Update(form *Form) error {
	//creating the query
	query := `
		UPDATE from
		SET user_id = $1, form_id = $2, form_status = $3, archive_status = $4, full_name = $5, other_names = $6, name_change_status = $7,
		social_security_num = $8, social_security_date = $9, social_security_country = $10, passport_number = $11,
		passport_date = $12, passport_country = $13, dob = $14, place_of_birth = $15, nationality = $16, acquired_nationality = $17,
		spouse_name = $18, address = $19, residential_phone_number = $20, residential_fax_number = $21, residential_email = $22, 
		version = version + 1
		WHERE form_id = $23
		AND version = $24
		RETURNING version
	`
	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//collecting the data fields into a slice for query insertion
	args := []interface{}{
		form.PublicUser_ID, form.ID, form.Status, form.Archive, form.Fullname, form.Othernames, form.Has_Changed_Name,
		form.SocialSecurity_Number, form.SocialSecurity_Date, form.SocialSecurity_Country, form.Passport_Number, form.Passport_Date,
		form.Passport_Country, form.DOB, form.Place_of_Birth, form.Nationality, form.Acquired_Nationality, form.Spouse_Name, form.Address,
		form.Phone_Number, form.Fax_Number, form.Residential_Email, form.ID, form.Version,
	}

	//Checke for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&form.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() - removes a specific form
func (m FormModel) Delete(id int64) error {
	//Ensuring that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}

	//creating the delete query
	query := `
		DELETE FROM form
		WHERE id = $1
	`

	//creating the context for this function
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Executing the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//Checking how many rows were affect by the delete operation, we call the RowsAffected() nethod on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//checking if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

//get all forms will be written later
