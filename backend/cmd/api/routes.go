// Filename: BIOAFF/backend/cmd/api
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	//httprouter instance and paths for handler functions
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/form", app.submitFormHandler)

	//not all middleware has been written yet only those that are needed to get the webserver runnning
	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
