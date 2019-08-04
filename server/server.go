package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = ":8000"

func Start() {
	http.HandleFunc("/", root)
	http.HandleFunc("/users/", Users)

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// https://stackoverflow.com/questions/40266633/golang-insert-null-into-sql-instead-of-empty-string
func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// $ curl http://localhost:8000/ -v
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Hello, World!")
})
