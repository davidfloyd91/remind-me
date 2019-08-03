package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/davidfloyd91/remind-me/db"
	"github.com/davidfloyd91/remind-me/types"
)

const port = ":8000"

func Start() {
	http.HandleFunc("/", root)
	http.HandleFunc("/users/", users)

	fmt.Println("Listening at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello, World!\n")
})

var users = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// $ curl http://localhost:8000/users/
	case "GET":
		var users []types.User

		query := `
        SELECT * FROM users
      `

		rows, err := db.DB.Query(query)
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var Id uint
			var Username string
			var Email string
			var Password string
			var CreatedAt time.Time
			var UpdatedAt time.Time
			var DeletedAt time.Time

			rows.Scan(&Id, &Username, &Email, &Password, &CreatedAt, &UpdatedAt, &DeletedAt)

			users = append(users, types.User{
				Id:        Id,
				Username:  Username,
				Email:     Email,
				Password:  Password,
				CreatedAt: CreatedAt,
				UpdatedAt: UpdatedAt,
				DeletedAt: DeletedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	// $ curl http://localhost:8000/users/ -d '{"Username":"Alice","Email":"alice@alice.alice", "Password":"lol"}' -v
	case "POST":
		user := &types.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var Id uint
		var Username string
		var Email string
		var CreatedAt time.Time

		query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at
      `

		err = db.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&Id, &Username, &Email, &CreatedAt)
		if err != nil {
			panic(err)
		}

		returnUser := types.User{
			Id:        Id,
			Username:  Username,
			Email:     Email,
			CreatedAt: CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returnUser)
	}
})
