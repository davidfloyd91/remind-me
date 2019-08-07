package server

import (
	"database/sql"
	"encoding/json"
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

	requestPath := r.URL.Path
	pathSplit := strings.Split(requestPath, "/")
	paramId := pathSplit[2]

	// for /users/:user_id/events/:event_id
	// send to server/handleEvents.go
	if len(pathSplit) > 3 {
		if pathSplit[3] == "events" {
			Events(w, r)
			return
		}
	}

	switch r.Method {
	case "GET":
		// $ curl http://localhost:8000/users/ -v | json_pp --json_opt=canonical,pretty
		// if paramId == "" {
		// 	query := `
    //       SELECT id, username, email, created_at, updated_at, deleted_at
    //       FROM users
    //       WHERE deleted_at IS NULL
    //       ORDER BY id
    //     `
		//
		// 	var err error
		// 	rows, err = db.DB.Query(query)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// } else {
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
		// }

	// $ curl http://localhost:8000/signup/ -d '{"Username":"Alice","Email":"alice@alice.alice", "Password":"lol"}' -v
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
      `

		_, err = db.DB.Query(query, user.Username, newNullString(user.Email), digest)
		if err != nil {
			// handle conflict error here
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString, err := user.GenerateJwt()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		token := types.Token{
			Token: tokenString,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&token)
		return

	// curl -X PATCH http://localhost:8000/users/5 -H "Token: lol" -d '{ "Email":"jaljdlkjadfj.email.email"}' -v
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
          password = COALESCE($3, password),
          updated_at = current_timestamp
        WHERE id = $4
        RETURNING id, username, email, created_at, updated_at, deleted_at
      `

		rows, err = db.DB.Query(query, newNullString(user.Username), newNullString(user.Email), newNullString(digest), paramId)
		if err != nil {
			// duplicate username fails silently here
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
}) // close Users

// $ curl http://localhost:8000/login/ -d '{"Username":"Egh","Password":"Superhilar"}' -v
var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var users []types.User
	user := &types.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `
    SELECT id, password
    FROM users
    WHERE username = $1
    AND deleted_at IS NULL
  `

	rows, err := db.DB.Query(query, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var Id uint
		var Password string

		rows.Scan(&Id, &Password)

		users = append(users, types.User{
			Id:       Id,
			Password: Password,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if len(users) == 0 {
		// right error code?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(user.Password))
	if err == nil {
		tokenString, err := users[0].GenerateJwt()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := types.Token{
			Token: tokenString,
		}
		json.NewEncoder(w).Encode(&token)
	} else {
		panic(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}) // close Login
