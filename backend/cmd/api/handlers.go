// Filename: BIOAFF/backend/cmd/api/handlers.go
package main

import (
	"fmt"
	"net/http"

	"se_group1.net/internal/data"
)

// submitFormHandler() - tries to create a form based on infromation supplied
func (app *application) submitFormHandler(w http.ResponseWriter, r *http.Request) {
	//Our target decode destination
	var input struct {
		PublicUser_ID          int64  `json:"puid"`
		Status                 string `json:"status"`
		Archive                bool   `json:"archive"`
		Fullname               string `json:"fullname"`
		Othernames             string `json:"othername"`
		Has_Changed_Name       bool   `json:"changed_name"`
		SocialSecurity_Number  int    `json:"ssnumber"`
		SocialSecurity_Date    string `json:"ssdate"`
		SocialSecurity_Country string `json:"sscountry"`
		Passport_Number        string `json:"passport_number"`
		Passport_Date          string `json:"passport_date"`
		Passport_Country       string `json:"passport_country"`
		DOB                    string `json:"dob"`
		Place_of_Birth         string `json:"place_of_birth"`
		Nationality            string `json:"nationality"`
		Acquired_Nationality   string `json:"acquired_nationality"`
		Spouse_Name            string `json:"spouse_name"`
		Address                string `json:"address"`
		Phone_Number           string `json:"phone"`
		Fax_Number             string `json:"fax"`
		Residential_Email      string `json:"residential_rmail"`
	}

	//initializing a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copying over the values from input to the new form
	form := &data.Form{
		PublicUser_ID:          input.PublicUser_ID,
		Status:                 input.Status,
		Archive:                input.Archive,
		Fullname:               input.Fullname,
		Othernames:             input.Othernames,
		Has_Changed_Name:       input.Has_Changed_Name,
		SocialSecurity_Number:  input.SocialSecurity_Number,
		SocialSecurity_Date:    input.SocialSecurity_Date,
		SocialSecurity_Country: input.SocialSecurity_Country,
		Passport_Number:        input.Passport_Number,
		Passport_Date:          input.Passport_Date,
		Passport_Country:       input.Passport_Country,
		DOB:                    input.DOB,
		Place_of_Birth:         input.Place_of_Birth,
		Nationality:            input.Nationality,
		Acquired_Nationality:   input.Acquired_Nationality,
		Spouse_Name:            input.Spouse_Name,
		Address:                input.Address,
		Phone_Number:           input.Phone_Number,
		Fax_Number:             input.Fax_Number,
		Residential_Email:      input.Residential_Email,
	}

	//no validation for now

	//creating the form
	err = app.models.Forms.Insert(form)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//creating a location header for the newly created resource/form
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/form/%d", form.ID))

	//writing jSON response with 201 - Created status code with the body
	//being the form data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"form": form}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
