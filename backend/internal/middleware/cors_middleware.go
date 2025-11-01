package middleware

import (
    "net/http"

    "github.com/gorilla/handlers"
)

func CORSMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
    return handlers.CORS(
        handlers.AllowedOrigins([]string{allowedOrigins}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )
}
