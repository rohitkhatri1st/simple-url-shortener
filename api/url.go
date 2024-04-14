package api

import (
	"fmt"
	"net/http"
	"simple-url-shortener/model"

	"github.com/gorilla/mux"
)

func (a *API) redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := mux.Vars(r)["shortKey"]
	notFoundPage := "http://localhost:8001/pages/404"
	if shortKey == "" {
		http.Redirect(w, r, notFoundPage, http.StatusNotFound)
		// a.SendJsonResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	originalUrl, err := a.App.Url.GetOriginalUrlFromShortKey(shortKey)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, notFoundPage, http.StatusNotFound)
		// a.SendJsonResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func (a *API) shortenUrl(w http.ResponseWriter, r *http.Request) {
	var shortenUrlRequest *model.ShortenUrlRequest
	if err := a.DecodeJSONBody(r, &shortenUrlRequest); err != nil {
		a.SendJsonResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	if errs := a.Validator.Validate(shortenUrlRequest); errs != nil {
		a.SendJsonResponse(w, http.StatusBadRequest, nil, errs...)
		return
	}
	shortenUrlResp, err := a.App.Url.ShortenUrl(shortenUrlRequest)
	if err != nil {
		a.SendJsonResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	a.SendJsonResponse(w, http.StatusOK, shortenUrlResp)
}
