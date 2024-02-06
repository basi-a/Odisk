package main

import (
	"context"
	"log"
	"net/http"
	g "odisk/global"
	"odisk/initialize"
	"os"
	"os/signal"
	"time"
)

func main()  {
	startAndShutdownServer()
}

func init()  {
	initialize.Initialize()
}

// startAndShutdownServer starts and gracefully shuts down an HTTP server
func startAndShutdownServer()  {
	// Create a new HTTP server instance
	srv := &http.Server{
		Addr: g.Config.Server.Port,
		Handler: g.RouterEngine,
	}
	// Start the server in a new goroutine
	go func ()  {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("listen: %s\n", err)
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