package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/TakedaB/akademi-api/internal/service"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "token não fornecido"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": "token inválido"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*service.Claims)
			if !ok {
				http.Error(w, `{"error": "não autenticado"}`, http.StatusUnauthorized)
				return
			}

			allowed := false
			for _, role := range roles {
				if claims.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, `{"error": "acesso negado"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
