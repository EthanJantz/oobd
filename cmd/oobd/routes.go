package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	return router
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := map[string]string{
		"status": "available",
	}

	err := app.writeJSON(w, http.StatusOK, env)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
