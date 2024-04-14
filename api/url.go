package api

import (
	"net/http"
)

func (a *API) shortenUrl(w http.ResponseWriter, r *http.Request) {
	a.SendJsonResponse(w, http.StatusOK, true)
}
