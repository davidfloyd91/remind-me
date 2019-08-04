package server

import (
	"database/sql"
	"encoding/json"
  // "fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/davidfloyd91/remind-me/db"
	"github.com/davidfloyd91/remind-me/types"
)

var Users = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var users []types.User
	var rows *sql.Rows

	paramId := strings.Split(r.URL.String(), "/")[2]

	switch r.Method {
	case "GET":
		// $ curl http://localhost:8000/users/ -v | json_pp --json_opt=canonical,pretty
		if paramId == "" {
			query := `
          SELECT id, username, email, created_at, updated_at, deleted_at
          FROM users
          WHERE deleted_at IS NULL
          ORDER BY id
        `

			var err error
			rows, err = db.DB.Query(query)
			if err != nil {
				panic(err)
			}
		} else {
      // $ curl http://localhost:8000/users/1 -v
			query := `
          SELECT id, username, email, created_at, updated_at, deleted_at
          FROM users
          WHERE id = $1
          AND deleted_at IS NULL
          ORDER BY id
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
      return
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
      return
		}

		digest := string(bytes)

		query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at, updated_at, deleted_at
      `

		rows, err = db.DB.Query(query, user.Username, user.Email, digest)
		if err != nil {
      // handle conflict error here
      w.WriteHeader(http.StatusBadRequest)
      return
		}

	// curl -X PATCH http://localhost:8000/users/1 -d '{"Username":"Egh","Password":"Superhilar"}' -v
	case "PATCH":
		if paramId == "" {
      w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := &types.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
      return
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
      return
		}

		digest := string(bytes)

		query := `
        UPDATE users SET
          username = COALESCE($1, username),
          email = COALESCE($2, email),
          password = COALESCE($3, password)
        WHERE id = $4
        RETURNING id, username, email, created_at, updated_at, deleted_at
      `

		rows, err = db.DB.Query(query, newNullString(user.Username), newNullString(user.Email), newNullString(digest), paramId)
		if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
		}

	// nothing is keeping us from seeing deleted users so far
	// curl -X DELETE http://localhost:8000/users/2 -v
	case "DELETE":
		if paramId == "" {
      w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := `
        UPDATE users SET
          deleted_at = current_timestamp
        WHERE id = $1
        RETURNING id, username, email, created_at, updated_at, deleted_at
      `

		var err error
		rows, err = db.DB.Query(query, paramId)
		if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
		}
	} // close switch

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

var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})
