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
	"time"
)

func main() {
	startAndShutdownServer()
}

func init() {
	initialize.Initialize()
}

// startAndShutdownServer starts and gracefully shuts down an HTTP server
func startAndShutdownServer() {
	// Create a new HTTP server instance
	srv := &http.Server{
		Addr:    g.Config.Server.Port,
		Handler: g.RouterEngine.Handler(),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	}

	// Check if the certificate and private key files exist
	cert := g.Config.Server.Ssl.Cert
	privateKey := g.Config.Server.Ssl.PrivateKey
	// log.Println(cert)
	// log.Println(privateKey)
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		log.Fatalf("Certificate file does not exist: %s\n", cert)
	}
	if _, err := os.Stat(privateKey); os.IsNotExist(err) {
		log.Fatalf("Private key file does not exist: %s\n", privateKey)
	}

	// Start the server in a new goroutine
	go func() {
		log.Printf("Check server with: curl -k -I https://localhost%s/ping", g.Config.Server.Port)

		if err := srv.ListenAndServeTLS(cert, privateKey); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s\n", err)
		}

	}()
	// Create an os.Signal channel and start listening for interrupt signals
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// Create a context with a timeout for the graceful shutdown of the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Attempt to gracefully shut down the server, waiting for all handlers to finish, or until the timeout
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
