package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
)

type Event struct {
	gorm.Model

	// ID        uint
	UserID      uint      `gorm:"user_id" json:"user_id"`
	Name        string    `gorm:"name" json:"name"`
	Description string    `gorm:"description" json:"description"`
	Scheduled   time.Time `gorm:"scheduled" json:"scheduled"`
}

type User struct {
	gorm.Model

	// ID        uint
	Username string `gorm:"username;unique;not null" json:"username"`
	Email    string `gorm:"email" json:"email"`
	Password string `gorm:"password" json:"password"`
}

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

type JWT struct {
	Token string
}

var db *gorm.DB

func main() {

	host := "localhost"
	port := 5432
	user := "postgres"
	dbname := "remindme"
	sslmode := "disable"

	t := "host=%s port=%d user=%s dbname=%s sslmode=%s"

	connectionString := fmt.Sprintf(t, host, port, user, dbname, sslmode)

	db, err := gorm.Open("postgres", connectionString)
  if err != nil {
		fmt.Println("Error in postgres connection: ", err)
	}

	db.AutoMigrate(&Event{}, &User{})
	defer db.Close()

	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	initRouter()
}

var Home = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Welcome home!")
}

var GetUsers = func(w http.ResponseWriter, r *http.Request) {
	var users []User

	db.Find(&users)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&users)
}

var UserEventsHelper = func(user_id string, timeframe string) []Event {
	var events []Event
	if timeframe == "all" {
		db.Where("user_id = ?", user_id).Find(&events)

	} else if timeframe == "today" {
		db.Where("user_id = ? and scheduled between now()::date and now()::date + interval '1d'", user_id).Find(&events)

	} else if timeframe == "tomorrow" {
		db.Where("user_id = ? and scheduled between now()::date + interval '1d' and now()::date + interval '2d'", user_id).Find(&events)

	} else if timeframe == "week" {
		db.Where("user_id = ? and scheduled between now()::date and now()::date + interval '7d'", user_id).Find(&events)

		// } else if timeframe == "next_week" {    //
		//     db.Where("user_id = ? and scheduled between now()::date + interval '1 week' and now()::date + interval '2 week'", user_id).Find(&events)
	}

	return events
}

var GetUserEvents = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	events := UserEventsHelper(params["user_id"], "all")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&events)
}

var GetUserEventsToday = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	events := UserEventsHelper(params["user_id"], "today")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&events)
}

var GetUserEventsTomorrow = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	events := UserEventsHelper(params["user_id"], "tomorrow")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&events)
}

var GetUserEventsThisWeek = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	events := UserEventsHelper(params["user_id"], "week")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&events)
}

//
// var GetUserEventsNextWeek = func(w http.ResponseWriter, r *http.Request) {
//     params := mux.Vars(r)
//     events := UserEventsHelper(params["user_id"], "next_week")
//
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(&events)
// }

var Signup = func(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	user.Password = string(bytes)

	err = db.Create(&user).Error
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&token)
}

var Login = func(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	storedUser := &User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("username = ?", user.Username).Find(&storedUser)

	w.Header().Set("Content-Type", "application/json")

	if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)) == nil {
		token, err := storedUser.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(&token)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (user User) GenerateJWT() (JWT, error) {
	signing_key := []byte(os.Getenv("JWT_SECRET"))

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_string, err := token.SignedString(signing_key)
	return JWT{token_string}, err
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		params := mux.Vars(r)

		if requestPath == "/" || requestPath == "/signup" || requestPath == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Token")

		w.Header().Set("Content-Type", "application/json")

		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		tk := &Claims{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
    if err != nil || !token.Valid {
      w.WriteHeader(http.StatusForbidden)
      return
    }

		if params["user_id"] != "" {
			id, err := strconv.ParseUint(params["user_id"], 10, 0)
			if err != nil || id != uint64(tk.UserID) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

var GetEvents = func(w http.ResponseWriter, r *http.Request) {
	var events []Event

	// add error handling
	db.Find(&events)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&events)
}

var GetEvent = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event

	// add error handling
	db.First(&event, params["id"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&event)
}

var CreateEvent = func(w http.ResponseWriter, r *http.Request) {
	event := &Event{}

	err := json.NewDecoder(r.Body).Decode(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// add error handling
	db.Create(&event)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

var UpdateEvent = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	event := &Event{}
	updateEvent := &Event{}

	err := json.NewDecoder(r.Body).Decode(updateEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.First(&event, params["event_id"])

	db.Model(&event).Updates(Event{
		Name:        updateEvent.Name,
		Description: updateEvent.Description,
		Scheduled:   updateEvent.Scheduled,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateEvent)
}

var DeleteEvent = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event

	// add error handling
	db.First(&event, params["event_id"])
	// add error handling
	db.Delete(&event)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&event)
}

func initRouter() {

	router := mux.NewRouter()

	// $ curl -H "Token: nicetrygithub" http://localhost:8000 -v
	router.HandleFunc("/", Home)

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users -v
	router.HandleFunc("/users", GetUsers).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users/1/events -v
	router.HandleFunc("/users/{user_id}/events", GetUserEvents).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users/1/events/today -v
	router.HandleFunc("/users/{user_id}/events/today", GetUserEventsToday).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users/1/events/tomorrow -v
	router.HandleFunc("/users/{user_id}/events/tomorrow", GetUserEventsTomorrow).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users/1/events/week -v
	router.HandleFunc("/users/{user_id}/events/week", GetUserEventsThisWeek).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/users/1/events/next -v
	// router.HandleFunc("/users/{user_id}/events/next", GetUserEventsNextWeek).Methods("GET")

	// $ curl http://localhost:8000/signup -d '{"username":"bigbaddude","email":"fun@fun.fun","password":"goodstuff"}' -v
	router.HandleFunc("/signup", Signup).Methods("POST")

	// $ curl http://localhost:8000/login -d '{"username":"noob", "password":"geeeez"}' -v
	router.HandleFunc("/login", Login).Methods("POST")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/events -v
	router.HandleFunc("/events", GetEvents).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/events/2 -v
	router.HandleFunc("/events/{id}", GetEvent).Methods("GET")

	// $ curl -H "Token: nicetrygithub" http://localhost:8000/events -d '{"name":"gee whiz","description":"you bet","user_id":11,"scheduled":"2019-06-30T23:10:00.000000-04:00"}' -v
	router.HandleFunc("/events", CreateEvent).Methods("POST")

	// $ curl need to do
	router.HandleFunc("/users/{user_id}/events/{event_id}", UpdateEvent).Methods("PUT")

	// $ curl -H "Token: nicetrygithub" -X DELETE http://localhost:8000/users/5/events/37 -v
	router.HandleFunc("/users/{user_id}/events/{event_id}", DeleteEvent).Methods("DELETE")

	// jwt
	router.Use(JwtAuthentication)

	// cors config
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"chrome-extension://nodpecpogkkofipgdkcchnbecnicoggl"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		// figure out which headers to allow
		AllowedHeaders: []string{"Token", "Host", "User-Agent", "Accept", "Content-Length", "Content-Type"},
	})

	// mux + jwt + cors
	handler := c.Handler(router)

	// start server
	log.Fatal(http.ListenAndServe(":8000", handler))
}
