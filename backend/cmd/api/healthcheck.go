//BIOAFF/backend/cmd/api/healthcheck.go
//Note: this is a testing file only

package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWrtier, r *http.Requst) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	//converting map -> JSON object
	err := app.writeJSON(w, http.StatusOK, data, nil)

	//will print error if any
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
