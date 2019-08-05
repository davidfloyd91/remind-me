package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func newNullTime(t time.Time) NullTime {
	if t.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return NullTime{}
	}
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

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

// this I guess: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func Start() {
	http.Handle("/", root)
	http.Handle("/users/", jwtAuthentication(Users))
	http.Handle("/signup/", Users)
	http.Handle("/login/", Login)

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// $ curl http://localhost:8000/ -v
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello, World!")
})
