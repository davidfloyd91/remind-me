package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

const port = ":8000"

func Start() {
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
	        AllowedOrigins: []string{"chrome-extension://nodpecpogkkofipgdkcchnbecnicoggl"},
	        AllowCredentials: true,
	        AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	        // figure out which headers to allow
	        AllowedHeaders: []string{"Token", "Host", "User-Agent", "Accept", "Content-Length", "Content-Type"},
	    })
  handler := c.Handler(mux)

	mux.Handle("/", root)
	mux.Handle("/signup/", Users)
	mux.Handle("/login/", Login)

	// redirects to Events
	mux.Handle("/users/", jwtAuthentication(Users))

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, handler))
}

// $ curl http://localhost:8000/ -v
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello, World!")
})
