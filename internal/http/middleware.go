package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (s *Server) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "empty auth header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "invalid auth header")
			return
		}

		userInfo, err := s.tokenManager.Parse(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtx, userInfo)))
	})
}
