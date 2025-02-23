package inboundhttp

import (
	"net/http"
	"strings"
)

func (server HttpServer) RequireValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/public/") {
			authToken := r.Header.Get("Authorization")
			authToken = strings.ReplaceAll(authToken, "Bearer ", "")
			err := server.authService.VerifyToken(r.Context(), authToken)
			if err != nil {
				writeResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func writeResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}
