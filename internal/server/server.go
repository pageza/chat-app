package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/routes"
	"github.com/sirupsen/logrus"
)

// StartServer initializes the HTTP server and listens for incoming requests.
func StartServer() {
	// Initialize Redis client
	rdb := redis.GetRedisClient()

	// Create a new router
	r := mux.NewRouter()
	r.Use(middleware.RateLimitMiddleware)

	// Add your routes here
	routes.InitializeRoutes(r, rdb)

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: r,
	}

	// Create a context for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	// Goroutine to listen for the interrupt signal
	go func() {
		<-stop // Wait for a signal to stop
		logrus.Info("Received signal, shutting down server.")

		// Create a context with a 5-second timeout for the server to close
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) // <-- Changed here
		defer cancel()                                                         // This is important

		// Attempt to gracefully shutdown the server
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Fatalf("Server Shutdown Failed:%+v", err)
		}
		// Removed the extra cancel() call
	}()

	// Start the server
	logrus.Infof("Server is running on port: %s", config.ServerPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("listen: %s\n", err)
	}

	// Wait for the context to be done
	<-ctx.Done()

	logrus.Info("Server stopped")
}
