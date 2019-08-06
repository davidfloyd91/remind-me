package types

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id        uint
	Username  string
	Email     string
	Password  string
	Admin     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Event struct {
	Id          uint
	UserId      uint
	Name        string
	Description string
	Scheduled   string // hacky
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Claims struct {
	UserId uint
	Admin bool
	jwt.StandardClaims
}
