package api

import (
	"net/http"
	"simple-url-shortener/model"
)

func (a *API) shortenUrl(w http.ResponseWriter, r *http.Request) {
	var shortenUrlRequest *model.ShortenUrlRequest
	if err := a.DecodeJSONBody(r, &shortenUrlRequest); err != nil {
		a.SendJsonResponse(w, http.StatusBadRequest, nil, err)
		return
	}
	if errs := a.Validator.Validate(shortenUrlRequest); errs != nil {
		a.SendJsonResponse(w, http.StatusInternalServerError, nil, errs...)
		return
	}
	shortenUrlResp, err := a.App.Url.ShortenUrl(shortenUrlRequest)
	if err != nil {
		a.SendJsonResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	a.SendJsonResponse(w, http.StatusOK, shortenUrlResp)
}
