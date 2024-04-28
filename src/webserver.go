package webserver

import (
	"ampapi-stats-wrapper/src/middleware"
	"ampapi-stats-wrapper/src/stats"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// WebServer - The web server
type WebServer struct {
	settings *stats.Settings
}

// NewWebServer - Create a new API server
func NewWebServer(settings *stats.Settings) *WebServer {
	return &WebServer{settings}
}

// Setup - Setup the web server
func (s *WebServer) Setup() http.Handler {
	router := http.NewServeMux()
	router = stats.ApplyRoutes(router, s.settings)
	return middleware.RequestLoggerMiddleware(router)
}

// Run - Start the web server
func (s *WebServer) Run() error {
	server := http.Server{
		Addr:    s.settings.ADDRESS,
		Handler: s.Setup(),
	}

	if s.settings.USE_UDS {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Remove(s.settings.ADDRESS)
			os.Exit(1)
		}()

		if _, err := os.Stat(s.settings.ADDRESS); err == nil {
			log.Printf("Removing existing socket file %s", s.settings.ADDRESS)
			if err := os.Remove(s.settings.ADDRESS); err != nil {
				return err
			}
		}

		socket, err := net.Listen("unix", s.settings.ADDRESS)
		if err != nil {
			return err
		}

		log.Printf("WebServer listening on %s", s.settings.ADDRESS)
		return server.Serve(socket)
	} else {
		log.Printf("WebServer listening on %s", s.settings.ADDRESS)
		return server.ListenAndServe()
	}
}
