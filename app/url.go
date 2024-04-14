package app

import (
	"simple-url-shortener/server/storage"

	"github.com/rs/zerolog"
)

type Url interface {
}

type UrlImplOpts struct {
	App    *App
	Db     storage.InMemoryDb
	Logger *zerolog.Logger
}
type UrlImpl struct {
	App    *App
	Db     storage.InMemoryDb
	Logger *zerolog.Logger
}

// InitUrl returns new instance of url implementation
func InitUrl(opts *UrlImplOpts) Url {
	l := opts.App.Logger.With().Str("service", "url").Logger()
	ui := UrlImpl{
		App:    opts.App,
		Db:     opts.Db,
		Logger: &l,
	}
	return &ui
}
