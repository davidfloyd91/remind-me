package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidfloyd91/remind-me/db"
	"github.com/davidfloyd91/remind-me/types"
)

var Events = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var events []types.Event
	var rows *sql.Rows

	requestPath := r.URL.Path
	pathSplit := strings.Split(requestPath, "/")
	paramUserId := pathSplit[2]
	paramEventId := pathSplit[4]

	switch r.Method {
	case "GET":
		// $ curl http://localhost:8000/users/5/events/ -H "Token: lol" -v | json_pp --json_opt=canonical,pretty
		query := `
        SELECT id, user_id, name, description, scheduled, created_at, updated_at, deleted_at
        FROM events
        WHERE user_id = $1
        AND deleted_at IS NULL
      `

		var err error
		rows, err = db.DB.Query(query, paramUserId)
		if err != nil {
			panic(err)
		}

	// $ curl http://localhost:8000/users/6/events/ -H "Token: lol" -d '{"Name":"yup yeah absolutely", "Description":"hrrrm", "Scheduled":"1908-01-02T03:04:05Z"}' -v
	case "POST":
		event := &types.Event{}

		err := json.NewDecoder(r.Body).Decode(event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		scheduled, err := time.Parse(time.RFC3339, event.Scheduled)
		if err != nil {
			panic(err)
		}

		query := `
          INSERT INTO events (user_id, name, description, scheduled)
          VALUES ($1, $2, $3, $4)
          RETURNING id, user_id, name, description, scheduled, created_at, updated_at, deleted_at
        `

		rows, err = db.DB.Query(query, paramUserId, event.Name, event.Description, scheduled)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// $ curl -X PATCH http://localhost:8000/users/6/events/5 -H "Token: lol" -d '{"Scheduled":"2087-01-02T15:04:05+02:00"}' -v
	case "PATCH":
		if paramEventId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		event := &types.Event{}

		err := json.NewDecoder(r.Body).Decode(event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		scheduled, err := time.Parse(time.RFC3339, event.Scheduled)
		if err != nil {
			panic(err)
		}

		query := `
            UPDATE events SET
              name = COALESCE($1, name),
              description = COALESCE($2, description),
              scheduled = COALESCE($3, scheduled),
              updated_at = current_timestamp
            WHERE id = $4
            RETURNING id, user_id, name, description, scheduled, created_at, updated_at, deleted_at
          `

		rows, err = db.DB.Query(query, newNullString(event.Name), newNullString(event.Description), scheduled, paramEventId)
		if err != nil {
			panic(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	// curl -X DELETE http://localhost:8000/users/6/events/8 -H "Token: lol" -v
	case "DELETE":
		if paramEventId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := `
		      UPDATE events SET
		        deleted_at = current_timestamp
		      WHERE id = $1
		      RETURNING id, user_id, name, description, scheduled, created_at, updated_at, deleted_at
		    `

		var err error
		rows, err = db.DB.Query(query, paramEventId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} // close switch

	for rows.Next() {
		var Id uint
		var UserId uint
		var Name string
		var Description string
		var Scheduled string
		var CreatedAt time.Time
		var UpdatedAt time.Time
		var DeletedAt time.Time

		rows.Scan(&Id, &UserId, &Name, &Description, &Scheduled, &CreatedAt, &UpdatedAt, &DeletedAt)

		events = append(events, types.Event{
			Id:          Id,
			UserId:      UserId,
			Name:        Name,
			Description: Description,
			Scheduled:   Scheduled,
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt,
			DeletedAt:   DeletedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
})
