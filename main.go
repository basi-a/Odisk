package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	g "odisk/global"
	"odisk/initialize"
	"os"
	"os/signal"
)

var srv *Server

// Server struct encapsulates the http.Server and its management
type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	// Create the http.Server instance
	s.httpServer = &http.Server{
		Addr:    g.Config.Server.Port,
		Handler: g.RouterEngine.Handler(),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	// Check for certificate and private key files
	cert := g.Config.Server.Ssl.Cert
	privateKey := g.Config.Server.Ssl.PrivateKey

	if err := s.checkFileExists(cert); err != nil {
		return err
	}
	if err := s.checkFileExists(privateKey); err != nil {
		return err
	}

	// Start the server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServeTLS(cert, privateKey); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// checkFileExists checks if a file exists and logs an error if not
func (s *Server) checkFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("%s file does not exist: %s\n", path, err)
	}
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initialize.Initialize()

	srv = &Server{}
	if err := srv.Start(ctx); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// Handle shutdown
	shutdownHandler(ctx)

	log.Println("Server exited gracefully")
}

func shutdownHandler(ctx context.Context) {
	// Create a signal channel and start listening for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Received interrupt signal, shutting down...")

	g.Producer.Stop()

	db, _ := g.DB.DB()
	err := db.Close()
	if err != nil {
		log.Println(err)
	}
	for _, v := range g.Consumers {
		v.Stop()
	}
	// Gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
