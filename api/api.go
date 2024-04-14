package api

import (
	"simple-url-shortener/app"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type API struct {
	Router     *Router
	MainRouter *mux.Router
	Logger     *zerolog.Logger

	App *app.App
}

type Router struct {
	Root    *mux.Router
	APIRoot *mux.Router
}

type Options struct {
	MainRouter *mux.Router
	Logger     *zerolog.Logger
}

func NewAPI(opts *Options) *API {
	api := API{
		MainRouter: opts.MainRouter,
		Router:     &Router{},
		Logger:     opts.Logger,
	}

	api.setupRoutes()
	return &api
}

func (a *API) setupRoutes() {
	a.Router.Root = a.MainRouter

	a.Router.APIRoot = a.MainRouter.PathPrefix("/api").Subrouter()
	// Declare More types of routes if needed
	a.InitRoutes()
}
