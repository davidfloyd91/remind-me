package main

import (
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    // "os"

    "golang.org/x/crypto/bcrypt"
    // "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    "github.com/rs/cors"
    _ "github.com/lib/pq"
    _ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
)

type Event struct {
	gorm.Model

  UserID       uint   `gorm:"user_id" json:"user_id"`
	Name         string `gorm:"name" json:"name"`
	Description  string `gorm:"description" json:"description"`
}

type User struct {
	gorm.Model

  Username     string `gorm:"username" json:"username"`
	Email        string `gorm:"email" json:"email"`
  Digest       string `json:"-"`
}

type JWTToken struct {
	Token        string `json:"token"`
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

    db.AutoMigrate(&Event{}, &User{})
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
               did that
    ******************************/

    initRouter()
}

var Home = func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Welcome home!")
}

var HashPassword = func(password string) string {
    bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
    return string(bytes)
}

var Signup = func(w http.ResponseWriter, r *http.Request) {
    user := &User{}

    // secure?
    err = json.NewDecoder(r.Body).Decode(user)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }

    user.Digest = HashPassword(r.FormValue("password"))

    db.Create(&user)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(&user)
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

    db.Model(&event).Updates(Event{Name: r.FormValue("name"), Description: r.FormValue("description")})

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

    /******************************
                 users
    ******************************/

    /*** signup ***/
    // $ curl -X POST http://localhost:8000/signup -d '{"username":"noob","email":"fun@fun.fun","password":"geeeez"}'
    router.HandleFunc("/signup", Signup).Methods("POST")

    /*** login ***/
    // $ curl
    // router.HandleFunc("/login", Login).Methods("POST")

    /******************************
                events
    ******************************/

    /*** index ***/
    // $ curl http://localhost:8000/events -v
    router.HandleFunc("/events", GetEvents).Methods("GET")

    /*** show ***/
    // $ curl http://localhost:8000/events/8 -v
    router.HandleFunc("/events/{id}", GetEvent).Methods("GET")

    /*** create ***/
    // $ curl -H "Content-Type: application/json" http://localhost:8000/events -d '{"name":"lol","description":"lol","user_id":1}' -v
    router.HandleFunc("/events", CreateEvent).Methods("POST")

    /*** update ***/
    // $ curl -X PUT http://localhost:8000/events/2 -d description=older -d name=so%20old -v
    router.HandleFunc("/events/{id}", UpdateEvent).Methods("PUT")

    /*** delete ***/
    // $ curl -X DELETE http://localhost:8000/events/10 -v
    router.HandleFunc("/events/{id}", DeleteEvent).Methods("DELETE")

    /******************************
                all done
    ******************************/

    // router.Use(JwtAuthentication)

    handler := cors.Default().Handler(router)

    log.Fatal(http.ListenAndServe(":8000", handler))
}
