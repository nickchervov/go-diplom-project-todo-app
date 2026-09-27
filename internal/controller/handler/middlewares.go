package handler

import (
	"net/http"
)

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.svc.IsAuthEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		var jwtToken string
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No cookie with key token", http.StatusUnauthorized)
			return
		}
		jwtToken = cookie.Value

		if err := h.svc.ValidateToken(jwtToken); err != nil {
			http.Error(w, "Authorization require", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
