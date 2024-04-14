package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"simple-url-shortener/api"
	"simple-url-shortener/app"
	"simple-url-shortener/app/validator"
	"simple-url-shortener/model"
	"simple-url-shortener/server/storage"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Server struct {
	httpServer *http.Server
	Router     *mux.Router
	Log        *zerolog.Logger
	API        *api.API
}

// NewServer returns a new Server object
func NewServer() *Server {
	r := mux.NewRouter()
	server := &Server{
		httpServer: &http.Server{},
		Router:     r,
	}

	server.InitLoggers()
	appLogger := server.Log.With().Str("type", "app").Logger()
	apiLogger := server.Log.With().Str("type", "api").Logger()

	// We can add a Middlewares here if required.
	// r.Use(authMiddleware)

	server.API = api.NewAPI(&api.Options{
		MainRouter: r,
		Logger:     &apiLogger,
		Validator:  validator.NewValidation(),
	})

	server.API.App = app.NewApp(&app.Options{Logger: &appLogger, Db: storage.InitDb()})

	app.InitService(server.API.App)
	return server
}

func (s *Server) InitLoggers() {
	cw := zerolog.ConsoleWriter{Out: os.Stdout}

	zlog := zerolog.New(cw).With().Timestamp().Stack().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	s.Log = &zlog
}

func (s *Server) StartServer() {
	addr := fmt.Sprintf("%s:%s", model.ServerHost, model.ServerPort)
	s.httpServer = &http.Server{
		Handler:      s.Router,
		Addr:         addr,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	s.Log.Info().Msgf("Staring server at %s", addr)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.Log.Error().Err(err).Msg("")
			return
		}
	}()
}

func (s *Server) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.Log.Debug().Msg("Shutting Down Server")
	s.httpServer.Shutdown(ctx)
	s.Log.Debug().Msg("HTTP Server Shut Down")
}
