package server

import (
  // "net/http"
	// "os"
  //
	// "github.com/dgrijalva/jwt-go"
)

// this is all wrong
// var jwtAuthentication = func(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         requestPath := r.URL.Path
//         params := mux.Vars(r)
//
//         if requestPath == "/" || requestPath == "/signup" || requestPath == "/login" {
//             next.ServeHTTP(w, r)
//             return
//         }
//
//         tokenHeader := r.Header.Get("Token")
//
//         w.Header().Set("Content-Type", "application/json")
//
//         if tokenHeader == "" {
//             w.WriteHeader(http.StatusForbidden)
//             return
//         }
//
//         tk := &Claims{}
//
//         token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
//             return []byte(os.Getenv("JWT_SECRET")), nil
//         })
//
//         if params["user_id"] != "" {
//             id, err := strconv.ParseUint(params["user_id"], 10, 0)
//             if err != nil || id != uint64(tk.UserId) {
//                 w.WriteHeader(http.StatusForbidden)
//                 return
//             }
//         }
//
//         if err != nil || !token.Valid {
//             w.WriteHeader(http.StatusForbidden)
//             return
//         }
//
//         next.ServeHTTP(w, r)
//     })
// }
