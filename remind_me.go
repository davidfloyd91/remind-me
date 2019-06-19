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
// right? wrong?
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

    // err????????
    db, err = gorm.Open("postgres", connectionString)

    db.AutoMigrate(&Event{})
    defer db.Close()

    // err????????
    if err != nil {
        fmt.Println("Error in postgres connection: ", err)
    }

    // err????????
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
    // error handling here (not all that important)
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Welcome home!")
}

var GetEvents = func(w http.ResponseWriter, r *http.Request) {
    var events []Event
    db.Find(&events)
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&events)
}

var CreateEvent = func(w http.ResponseWriter, r *http.Request) {
    event := &Event{}

    // should it be declaring a new err?
    err := json.NewDecoder(r.Body).Decode(event)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    db.Create(&event)

    // error handling here
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(event)
}

func initRouter() {
    router := mux.NewRouter()

    router.HandleFunc("/", Home)
    router.HandleFunc("/events", GetEvents).Methods("GET")
    router.HandleFunc("/events", CreateEvent).Methods("POST")

    // router.Use(JwtAuthentication)

    handler := cors.Default().Handler(router)

    // declare new err here?
    err := http.ListenAndServe(":8000", handler)

    // which err is this even referring to?
	  if err != nil {
		    fmt.Print(err)
    }
}
