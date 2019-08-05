package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// this is behaving very strangely
// func newNullTime(t time.Time) pq.NullTime {
// 	if t.String() == "0001-01-01 00:00:00 +0000 UTC" {
// 		return pq.NullTime{}
// 	}
// 	return pq.NullTime{
// 		Time:  t,
// 		Valid: true,
// 	}
// }

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

// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func Start() {
	http.Handle("/", root)
	http.Handle("/signup/", Users)
	http.Handle("/login/", Login)

	// redirects to Events
	http.Handle("/users/", jwtAuthentication(Users))

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// $ curl http://localhost:8000/ -v
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello, World!")
})
