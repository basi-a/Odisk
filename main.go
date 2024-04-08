// package main

// import (
// 	"context"
// 	"crypto/tls"
// 	"log"
// 	"net/http"
// 	g "odisk/global"
// 	"odisk/initialize"
// 	"os"
// 	"os/signal"
// 	"time"
// )

// var Srv *http.Server

// func main() {
// 	startAndShutdownServer()
// 	quit := make(chan os.Signal, 10)
// 	signal.Notify(quit, os.Interrupt)
// 	go QuitSignalHandler(quit)
// }

// func init() {
// 	initialize.Initialize()
// }

// // startAndShutdownServer starts and gracefully shuts down an HTTP server
// func startAndShutdownServer() {
// 	// Create a new HTTP server instance
// 	Srv = &http.Server{
// 		Addr:    g.Config.Server.Port,
// 		Handler: g.RouterEngine.Handler(),
// 		TLSConfig: &tls.Config{
// 			MinVersion:               tls.VersionTLS12,
// 			PreferServerCipherSuites: true,
// 		},
// 	}

// 	// Check if the certificate and private key files exist
// 	cert := g.Config.Server.Ssl.Cert
// 	privateKey := g.Config.Server.Ssl.PrivateKey

// 	if _, err := os.Stat(cert); os.IsNotExist(err) {
// 		log.Fatalf("Certificate file does not exist: %s\n", cert)
// 	}
// 	if _, err := os.Stat(privateKey); os.IsNotExist(err) {
// 		log.Fatalf("Private key file does not exist: %s\n", privateKey)
// 	}

// 	// Start the server in a new goroutine
// 	go func() {
// 		log.Printf("Check server with: curl -k -I https://localhost%s/ping", g.Config.Server.Port)

// 		if err := Srv.ListenAndServeTLS(cert, privateKey); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Failed to start server: %s\n", err)
// 		}

// 	}()
// 	// // Create an os.Signal channel and start listening for interrupt signals
// 	// quit := make(chan os.Signal, 10)
// 	// signal.Notify(quit, os.Interrupt)
// 	// <-quit
// 	// log.Println("Shutdown Server ...")

// 	// // Create a context with a timeout for the graceful shutdown of the server
// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()
// 	// // Attempt to gracefully shut down the server, waiting for all handlers to finish, or until the timeout
// 	// if err := Srv.Shutdown(ctx); err != nil {
// 	// 	log.Fatal("Server Shutdown:", err)
// 	// }
// 	// log.Println("Server exiting")
// }

// func QuitSignalHandler(quit <-chan os.Signal) {
// 	log.Println("Received signal:", <-quit)
// 	g.Producer.Stop()
// 	g.Consumer.Stop()
// 	log.Println("Shutdown Server ...")

// 	// Create a context with a timeout for the graceful shutdown of the server
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	// Attempt to gracefully shut down the server, waiting for all handlers to finish, or until the timeout
// 	if err := Srv.Shutdown(ctx); err != nil {
// 		log.Fatal("Server Shutdown:", err)
// 	}
// 	log.Println("Server exited")

// }
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
    g.Consumer.Stop()

    // Gracefully shut down the server
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
}