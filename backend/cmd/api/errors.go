// Filename: BIOAFF/backend/cmd/api/errors.go
package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// we want to send JSON-formatted error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//Create the JSON response
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Server error response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//we log the error
	app.logError(r, err)

	//Prepare a message with the error
	message := "the server encounted a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// /The not found response
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//Create our message
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// A method not allowed response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	//Create our message
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// User provide a bad request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	//Create our message
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Validation error
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// Edit conflict error
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Rate limit error
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// Invalid credentials
func (app *application) invalidCredntialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Invalid Token
func (app *application) invalidAuthenticationTokenReponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("WWW-Authenicate", "Bearer")
	message := "invalid or missing authorization token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Unauthorized access
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Users who have not activated their account
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// User does not have the required permission
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user accound does not have the necessary permission to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
