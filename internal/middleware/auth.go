package middleware

import (
	"context"
	"net/http"
	"strings"

	"vatsim-auth-service/internal/jwt"
)

type ContextKey string

const (
	ContextKeyCID   ContextKey = "cid"
	ContextKeyEmail ContextKey = "email"
	ContextKeyRoles ContextKey = "roles"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		header := r.Header.Get("Authorization")
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		} else {
			// Проверка куки
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				tokenStr = cookie.Value
			}
		}

		if tokenStr == "" {
			http.Error(w, "missing or malformed Authorization header or auth_token cookie", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyCID, claims.CID)
		ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextKeyRoles, claims.Roles)
		ctx = context.WithValue(ctx, ContextKey("country_name"), claims.CountryName)
		ctx = context.WithValue(ctx, ContextKey("division_name"), claims.DivisionName)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
