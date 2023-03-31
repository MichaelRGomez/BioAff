// Filename: BIOAFF/backend/cmd/api/healtCheck.go
package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//creating map that'll hold the health data
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	//coverting map -> JSON object
	err := app.writeJSON(w, http.StatusOK, data, nil)

	//will print error if any
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
