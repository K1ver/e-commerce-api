package middleware

import (
	"net/http"
	"strings"

	jwtmanager "github.com/K1ver/e-commerce-api/internal/pkg/jwt"
	"github.com/K1ver/e-commerce-api/internal/pkg/ctxkey"
)

type Auth struct {
	jwt *jwtmanager.Manager
}

func NewAuth(jwt *jwtmanager.Manager) *Auth {
	return &Auth{jwt: jwt}
}

func (m *Auth) RequireJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		subject, err := m.jwt.ParseAccess(token)
		if err != nil {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := ctxkey.WithAuth(r.Context(), subject.UserID, subject.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}
