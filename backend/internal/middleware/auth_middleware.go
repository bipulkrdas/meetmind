package middleware

import (
	"context"
	"net/http"
	"strings"

	"livekit-consulting/backend/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtSecret string, userRepo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Could not find bearer token in Authorization header", http.StatusUnauthorized)
				return
			}

			claims := &jwt.StandardClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := WithUser(r.Context(), user)
			ctx = context.WithValue(ctx, "userID", claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
