package model

type ShortenUrlRequest struct {
	OriginalUrl string `json:"original_url" validate:"required,url"`
}

type UrlData struct {
	OriginalUrl string `json:"original_url"`
	ShortenUrl  string `json:"shorten_url"`
	ShortKey    string `json:"shortKey"`
}

type ShortenUrlResp struct {
	OriginalUrl string `json:"original_url"`
	ShortenUrl  string `json:"shorten_url"`
}

type DomainCount struct {
	Domain string
	Count  int
}

const ShortKeyLength = 6

const (
	ServerScheme = "http"
	ServerHost   = "0.0.0.0"
	ServerPort   = "8001"
)
