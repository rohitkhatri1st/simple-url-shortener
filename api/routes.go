package api

import (
	"net/http"
)

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	// We can create a handler function separately that may act as a middleware for additional functionalities.
	a.Router.APIRoot.HandleFunc("/health-check", a.healthCheck).Methods(http.MethodGet)

	a.Router.APIRoot.HandleFunc("/shorten", a.shortenUrl).Methods(http.MethodPost)
	a.Router.Root.HandleFunc("/{shortKey}", a.redirect).Methods(http.MethodGet)
}

func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	a.SendJsonResponse(w, http.StatusOK, true)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// jsonResponse, _ := json.Marshal(true)
	// w.Write(jsonResponse)
}
