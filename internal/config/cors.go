// config/cors_config.go

package config

import (
	"github.com/rs/cors"
)

// InitializeCORS initializes the CORS middleware and returns it
func InitializeCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   CorsAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   CorsAllowedMethods,
		AllowedHeaders:   CorsAllowedHeaders,
	})
}
