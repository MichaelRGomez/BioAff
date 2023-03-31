// Filename: BIOAFF/backend/cmd/api/main_text.go
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (app *application) testsubmitFormHandler(t *testing.T) {
	//json request
	body := []byte(`{
		"puid":"1",
		"status":"pending",
		"archive":false,
		"fullname":"jd",
		"othername":"",
		"changed_name":false,
		"ssnumber":"000000000",
		"ssdate":"1990-02-02",
		"sscountry":"belize",
		"passport_number":"AD142641",
		"passport_date":"2001-04-01",
		"passport_country":"belize",
		"dob":"1980-07-14",
		"place_of_birth":"KHMH",
		"nationality":"belizean",
		"acquired_nationality":"birth",
		"spouse_name":"",
		"address":"17 starstoon street",
		"phone":"5016045434",
		"fax":"",
		"residential_rmail":""
	}`)
	request, err := http.NewRequest("POST", "/v1/form", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.submitFormHandler)
	handler.ServeHTTP(rr, request)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"form_id":1, "created_at":"", "version":1}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
