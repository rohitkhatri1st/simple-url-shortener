package api

import (
	"encoding/json"
	"net/http"
)

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.HandleFunc("/health-check", a.healthCheck).Methods("GET")
}

func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(true)
	w.Write(jsonResponse)
}
