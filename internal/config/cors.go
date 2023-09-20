// Package config contains configuration settings and initializers for the chat application.
// This file specifically deals with the CORS (Cross-Origin Resource Sharing) configuration.

package config

import (
	"github.com/rs/cors"
)

// InitializeCORS sets up and returns the CORS middleware.
// It uses the configuration settings defined in the config package.
func InitializeCORS() *cors.Cors {
	// Create a new CORS middleware with specific options
	return cors.New(cors.Options{
		AllowedOrigins:   CorsAllowedOrigins,        // Only allow specific origins to access resources
		AllowCredentials: true,                      // Allow cookies and authentication headers
		AllowedMethods:   CorsAllowedMethods,        // Only allow specific HTTP methods (e.g., GET, POST)
		AllowedHeaders:   CorsAllowedHeaders,        // Only allow specific HTTP headers
		ExposedHeaders:   []string{"Authorization"}, // Explicitly expose the Authorization header to clients
		MaxAge:           600,                       // Cache CORS preflight requests for 10 minutes
	})
}
