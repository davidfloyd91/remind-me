package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
    Id int `json: "id", db: "id"`
    Username string `json:"username", db:"username"`
    Password string `json:"password", db:"password"`
}

type Event struct {
    Id int `json: "id", db: "id"`
    UserId int `json:"user_id", db:"user_id"`
    Name string `json:"name", db:"name"`
    Description string `json:"description", db:"description"`
    Time time.Time `json:"time", db:"time"`
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/", Home)

    router.HandleFunc("/login", Login)
    router.HandleFunc("/signup", Signup)

    // http.HandleFunc("/events/", RetrieveEvents)
    router.HandleFunc("/events/new", CreateEvent)
    // http.HandleFunc("/events/{id}/edit", EditEvent)

    initDB()

    log.Fatal(http.ListenAndServe(":8000", router))
}

func Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Sup World!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	creds := &User{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
  		w.WriteHeader(http.StatusBadRequest)
  		return
	}

	result := db.QueryRow("select password from users where username=$1", creds.Username)
	if err != nil {
  		w.WriteHeader(http.StatusInternalServerError)
  		return
	}

	storedCreds := &User{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
  		if err == sql.ErrNoRows {
    			w.WriteHeader(http.StatusUnauthorized)
    			return
  		}

  		w.WriteHeader(http.StatusInternalServerError)
  		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		  w.WriteHeader(http.StatusUnauthorized)
	}
}

// curl -X POST -H 'Content-Type:application/json' http://localhost:8000/signup -d '{"username":"blamo","password":"blamo"}' -v

// curl -X POST -H 'Content-Type:application/json' http://localhost:8000/login -d '{"username":"blamo","password":"blamo"}' -v

// curl -X POST -H 'Content-Type:application/json' http://localhost:8000/events/new -d '{"user_id":4,"name":"wow","description":"whooooa"}' -v

func Signup(w http.ResponseWriter, r *http.Request) {
    creds := &User{}
    err := json.NewDecoder(r.Body).Decode(creds)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

    if _, err = db.Query("insert into users (username, password) values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
    event := &Event{}
    err := json.NewDecoder(r.Body).Decode(event)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if _, err = db.Query("insert into events (user_id, name, description) values ($1, $2, $3)", event.UserId, event.Name, event.Description); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func initDB() {
    var err error
    db, err = sql.Open("postgres", "dbname=remindme sslmode=disable")

    if err != nil {
        panic(err)
    }
}

// https://github.com/sohamkamani/go-password-auth-example
// https://golang.org/doc/articles/wiki/
// https://tutorialedge.net/golang/creating-restful-api-with-golang/
// https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
