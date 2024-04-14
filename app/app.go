package app

import (
	"simple-url-shortener/server/storage"

	"github.com/rs/zerolog"
)

type App struct {
	Logger *zerolog.Logger
	Db     *storage.InMemoryDatabases
	// List of services this app is implementing
	Url Url
}

// Options contains arguments required to create a new app instance
type Options struct {
	Logger *zerolog.Logger
	Db     *storage.InMemoryDatabases
}

func NewApp(opts *Options) *App {
	return &App{
		Logger: opts.Logger,
		Db:     opts.Db,
	}
}
