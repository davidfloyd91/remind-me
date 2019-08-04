package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type NullString struct {
  String string
  Valid  bool // Valid is true if String is not NULL
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func newNullString(s string) NullString {
	if len(s) == 0 {
		return NullString{}
	}
	return NullString{
		String: s,
		Valid:  true,
	}
}

func newNullTime(t time.Time) NullTime {
	if t.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return NullTime{}
	}
	return NullTime{
		Time: t,
		Valid: true,
	}
}

const port = ":8000"

func Start() {
	http.HandleFunc("/", root)
	http.HandleFunc("/users/", Users)

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// $ curl http://localhost:8000/ -v | json_pp
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello, World!")
})
