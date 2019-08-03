package server

import (
  "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "strings"
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

// $ curl http://localhost:8000/ -v
var root = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode("Hello, World!")
})

var users = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  var users []types.User
  var rows *sql.Rows

	switch r.Method {
	case "GET":
    paramId := strings.Split(r.URL.String(), "/")[2]

    // $ curl http://localhost:8000/users/ -v
    if paramId == "" {
  		query := `
          SELECT id, username, email, created_at, updated_at, deleted_at
          FROM users
        `

      var err error
  		rows, err = db.DB.Query(query)
  		if err != nil {
  			panic(err)
  		}

    // $ curl http://localhost:8000/users/3 -v
    } else {
      query := `
          SELECT id, username, email, created_at, updated_at, deleted_at
          FROM users
          WHERE id = $1
        `

      var err error
      rows, err = db.DB.Query(query, paramId)
      if err != nil {
  			panic(err)
  		}
    }

	// $ curl http://localhost:8000/users/ -d '{"Username":"Alice","Email":"alice@alice.alice", "Password":"lol"}' -v
	case "POST":
		user := &types.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

    // passwords are not being hashed -- must fix
		query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at, updated_at, deleted_at
      `

		rows, err = db.DB.Query(query, user.Username, user.Email, user.Password)
		if err != nil {
			panic(err)
		}
	}

  for rows.Next() {
    var Id uint
    var Username string
    var Email string
    var CreatedAt time.Time
    var UpdatedAt time.Time
    var DeletedAt time.Time

    rows.Scan(&Id, &Username, &Email, &CreatedAt, &UpdatedAt, &DeletedAt)

    users = append(users, types.User{
      Id:        Id,
      Username:  Username,
      Email:     Email,
      CreatedAt: CreatedAt,
      UpdatedAt: UpdatedAt,
      DeletedAt: DeletedAt,
    })
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(users)
})
