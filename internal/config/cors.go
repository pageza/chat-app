// config/cors_config.go

package config

import (
	"github.com/rs/cors"
)

// InitializeCORS initializes the CORS middleware and returns it
func InitializeCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   CorsAllowedOrigins, // Only allow specific origins
		AllowCredentials: true,
		AllowedMethods:   CorsAllowedMethods,        // Only allow specific methods (e.g., GET, POST)
		AllowedHeaders:   CorsAllowedHeaders,        // Only allow specific headers
		ExposedHeaders:   []string{"Authorization"}, // Expose Authorization header
		MaxAge:           600,                       // Cache preflight request for 10 minutes
	})
}
