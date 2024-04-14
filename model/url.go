package model

type ShortenUrlRequest struct {
	OriginalUrl string `json:"original_url" validate:"required"`
}

type UrlData struct {
	OriginalUrl string `json:"original_url"`
	ShortenUrl  string `json:"shorten_url"`
}

type ShortenUrlResp struct {
	OriginalUrl string `json:"original_url"`
	ShortenUrl  string `json:"shorten_url"`
}
