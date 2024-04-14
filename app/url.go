package app

import "github.com/rs/zerolog"

type Url interface {
}

type UrlImplOpts struct {
	App    *App
	Logger *zerolog.Logger
}
type UrlImpl struct {
	App    *App
	Logger *zerolog.Logger
}

// InitUrl returns new instance of url implementation
func InitUrl(opts *UrlImplOpts) Url {
	l := opts.App.Logger.With().Str("service", "url").Logger()
	ui := UrlImpl{
		App:    opts.App,
		Logger: &l,
	}
	return &ui
}
