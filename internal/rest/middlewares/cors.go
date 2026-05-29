package middlewares

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CorsDocs sets CORS headers for the SwaggerDoc UI to be able to access the API.
func (p *Provider) CorsDocs() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
