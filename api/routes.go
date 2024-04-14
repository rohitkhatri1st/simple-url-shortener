package api

import (
	"encoding/json"
	"net/http"
)

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	// We can create a handler function separately that may act as a middleware for additional functionalities.
	a.Router.Root.HandleFunc("/health-check", a.healthCheck).Methods(http.MethodGet)

	a.Router.Root.HandleFunc("/url", a.shortenUrl).Methods(http.MethodGet)
}

func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(true)
	w.Write(jsonResponse)
}
