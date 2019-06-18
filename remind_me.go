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
    initDB()
    initRouter()
}

var Home = func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Welcome home!")
}

// this doesn't work
var CreateEvent = func(w http.ResponseWriter, r *http.Request) {
    event := &Event{}
    json.NewDecoder(r.Body).Decode(&event)

    db.Create(&event)
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&event)
}

func initRouter() {
    router := mux.NewRouter()

    router.HandleFunc("/", Home)
    router.HandleFunc("/events", CreateEvent).Methods("POST")

    // router.Use(JwtAuthentication)

    handler := cors.Default().Handler(router)

    err := http.ListenAndServe(":8000", handler)
	  if err != nil {
		    fmt.Print(err)
    }
}

// https://github.com/Sach97/pgsearch-gorm-example/blob/master/search/main.go
func initDB() {
    host := "localhost"
    port := 5432
    user := "postgres"
    dbname := "remindme"
    sslmode := "disable"

    t := "host=%s port=%d user=%s dbname=%s sslmode=%s"

    connectionString := fmt.Sprintf(t, host, port, user, dbname, sslmode)

    db, err := gorm.Open("postgres", connectionString)

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
}
