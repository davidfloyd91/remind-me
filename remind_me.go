package main

import (
    "fmt"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    _ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
)

type Event struct {
	gorm.Model

	Name         string
	Description  string
}

var db *gorm.DB
var err error

func main() {
    /******************************
        initialize the database
    ******************************/
    host := "localhost"
    port := 5432
    user := "postgres"
    dbname := "remindme"
    sslmode := "disable"

    t := "host=%s port=%d user=%s dbname=%s sslmode=%s"

    connectionString := fmt.Sprintf(t, host, port, user, dbname, sslmode)

    db, err = gorm.Open("postgres", connectionString)

    db.AutoMigrate(&Event{})
    defer db.Close()

    if err != nil {
        fmt.Println("Error in postgres connection: ", err)
    }

    err = db.DB().Ping()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Successfully connected!")
    /******************************
        .............. did that
    ******************************/

    initRouter()
}

var Home = func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Welcome home!")
}

var GetEvent = func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var event Event

    db.First(&event, params["id"])

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&event)
}

var UpdateEvent = func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
  	var event Event

  	db.First(&event, params["id"])
    r.ParseForm()
    db.Model(&event).Updates(Event{Name: r.FormValue("Name"), Description: r.FormValue("Description")})

    w.Header().Set("Content-Type", "application/json")
  	json.NewEncoder(w).Encode(&event)
}

var DeleteEvent = func(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var event Event

    db.First(&event, params["id"])
    db.Delete(&event)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&event)
}

var GetEvents = func(w http.ResponseWriter, r *http.Request) {
    var events []Event

    db.Find(&events)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&events)
}

var CreateEvent = func(w http.ResponseWriter, r *http.Request) {
    event := &Event{}

    err = json.NewDecoder(r.Body).Decode(event)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    db.Create(&event)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(event)
}

func initRouter() {
    router := mux.NewRouter()

    /*** home ***/
    // $ curl http://localhost:8000 -v
    router.HandleFunc("/", Home)

    /*** index ***/
    // $ curl http://localhost:8000/events -v
    router.HandleFunc("/events", GetEvents).Methods("GET")

    /*** show ***/
    // $ curl http://localhost:8000/events/8 -v
    router.HandleFunc("/events/{id}", GetEvent).Methods("GET")

    /*** create ***/
    // $ curl -H "Content-Type: application/json" http://localhost:8000/events -d '{"name":"newnewnew","description":"whoa such novelty"}' -v
    router.HandleFunc("/events", CreateEvent).Methods("POST")

    /*** update ***/
    // $ curl -X PUT http://localhost:8000/events/10 -d Description=older -d Name=so%20old -v
    router.HandleFunc("/events/{id}", UpdateEvent).Methods("PUT")

    /*** delete ***/
    // $ curl -X DELETE http://localhost:8000/events/10 -v
    router.HandleFunc("/events/{id}", DeleteEvent).Methods("DELETE")

    // router.Use(JwtAuthentication)

    handler := cors.Default().Handler(router)

    log.Fatal(http.ListenAndServe(":8000", handler))
}
