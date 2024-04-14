package main

import (
	"os"
	"os/signal"
	"simple-url-shortener/server"
	"syscall"
)

func main() {
	s := server.NewServer()
	s.StartServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	s.StopServer()
	s.Log.Info().Msg("Server Stopped")
}
