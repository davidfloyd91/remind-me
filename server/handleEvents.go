package server

import (
  "database/sql"
  "encoding/json"
  // "fmt"
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
  paramId := strings.Split(requestPath, "/")[2]

  switch r.Method {
    // $ curl http://localhost:8000/users/5/events/ -H "Token: lol" -d '{ "Name":"A real good one", "Description":"Oh geez"}' -v
    case "POST":
      event := &types.Event{}

      err := json.NewDecoder(r.Body).Decode(event)
      if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
      }

      query := `
          INSERT INTO events (user_id, name, description, scheduled)
          VALUES ($1, $2, $3, $4)
          RETURNING id, user_id, name, description, scheduled, created_at, updated_at, deleted_at
        `

      rows, err = db.DB.Query(query, paramId, event.Name, event.Description, event.Scheduled)
      if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
      }
    }

    for rows.Next() {
      var Id uint
      var UserId uint
      var Name string
      var Description string
      var Scheduled time.Time
      var CreatedAt time.Time
      var UpdatedAt time.Time
      var DeletedAt time.Time

      rows.Scan(&Id, &UserId, &Name, &Description, &Scheduled, CreatedAt, &UpdatedAt, &DeletedAt)

      events = append(events, types.Event{
        Id:        Id,
  			UserId:    UserId,
  			Name:      Name,
        Description:      Description,
        Scheduled:      Scheduled,
  			CreatedAt: CreatedAt,
  			UpdatedAt: UpdatedAt,
  			DeletedAt: DeletedAt,
      })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)

})
