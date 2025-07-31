package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JsonError(w, http.StatusUnauthorized, "authorization header required")
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			JsonError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := parts[1]
		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			JsonError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Добавляем пользователя в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Вспомогательная функция для получения пользователя из контекста
func GetUserFromContext(ctx context.Context) *auth.Claims {
	if claims, ok := ctx.Value(UserContextKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}
