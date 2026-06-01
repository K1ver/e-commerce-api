package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowOrigins string) func(http.Handler) http.Handler {
	origins := strings.Split(allowOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && (allowOrigins == "*" || originAllowed(origins, origin)) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if allowOrigins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func originAllowed(origins []string, origin string) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}
