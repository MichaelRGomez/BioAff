// BIOAFF/backend/cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routers() http.Handler {
	//httprouter instance and the paths for each handler function

	//the router instance
	router := httprouter.New()

	//paths
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//return; with all middleware layered on
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
