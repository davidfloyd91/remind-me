package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/davidfloyd91/remind-me/types"
)

var jwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := &types.Claims{}
		requestPath := r.URL.Path
		pathSplit := strings.Split(requestPath, "/")
		paramId := pathSplit[2]

		w.Header().Set("Content-Type", "application/json")

		tokenHeader := r.Header.Get("Token")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if paramId != "" {
			id, err := strconv.ParseUint(paramId, 10, 0)
			if err != nil || id != uint64(tk.UserId) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
