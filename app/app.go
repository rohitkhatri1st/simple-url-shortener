package app

import "github.com/rs/zerolog"

type App struct {
	Logger *zerolog.Logger
	// List of services this app is implementing
	Url Url
}

// Options contains arguments required to create a new app instance
type Options struct {
	Logger *zerolog.Logger
}

func NewApp(opts *Options) *App {
	return &App{
		Logger: opts.Logger,
	}
}
