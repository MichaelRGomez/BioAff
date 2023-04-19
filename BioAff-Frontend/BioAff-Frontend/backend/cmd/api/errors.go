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

// errorResponse() - sends a JSON-formatted error message to the client
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//creating the json response
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse() - reports on errors that occur on the server
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//logging the error
	app.logError(r, err)

	//preparing a message with the error
	message := "the server encountered a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// 404 not found - for when the path isn't a registered path
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//creating our message
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// bad request - when the user supplies a badly formatted request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	//creating on error message
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Validation error - for when something goes wrong during validation
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// Edit conflict error - for when something goes wrong with editing a db record(s)
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Rate limit error - once something tries to pass the rate limit
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// Invalid credentials
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// Invalid token
func (app *application) invalidAuthenicationToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("WWW-Authenticate", "Bearer")
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

// User does not have required permission
func (app *application) notPermittedReponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account does not have the necessary permission to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
