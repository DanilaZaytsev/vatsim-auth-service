package middleware

import (
	"net/http"
)

// RequireRolesMiddleware проверяет, что у пользователя одна из допустимых ролей
func RequireRolesMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if _, ok := roleSet[claims.Roles]; !ok {
				http.Error(w, "forbidden — insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
