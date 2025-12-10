package middleware

import (
	"net/http"

	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
)

// RequireJSON middleware para validar que el Content-Type sea application/json
func RequireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			// Permitir multipart/form-data para uploads
			if contentType != "application/json" && contentType != "" {
				if contentType != "multipart/form-data" && !isMultipart(contentType) {
					utils.JSONError(w, http.StatusUnsupportedMediaType, "Content-Type debe ser application/json")
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isMultipart(contentType string) bool {
	return len(contentType) >= 19 && contentType[:19] == "multipart/form-data"
}
