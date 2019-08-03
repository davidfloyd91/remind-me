package types

import (
	"time"
)

type User struct {
	Id        uint
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Event struct {
	Id          uint
	UserId      uint
	Name        string
	Description string
	scheduled   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
