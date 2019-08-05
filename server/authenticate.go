package server

import (
  "fmt"
  "net/http"
	"os"
  "strconv"
  "strings"

	"github.com/dgrijalva/jwt-go"

  "github.com/davidfloyd91/remind-me/types"
)

// $ curl -X PATCH http://localhost:8000/users/5 -d '{"Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjUsImV4cCI6MTU2NzU2ODY5OX0.u4sIkIrpY9K6oWcjPTsAaPFKD_3ua9MKR7Ye0q3fSNE", "Email":"jaljdlkjadfj.email.email"}' -v
var jwtAuthentication = func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tk := &types.Claims{}
        requestPath := r.URL.Path
        paramId := strings.Split(requestPath, "/")[2]

        fmt.Println(paramId) // 5

        w.Header().Set("Content-Type", "application/json")

        tokenHeader := r.Header.Get("Token")
        fmt.Println(tokenHeader)
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
