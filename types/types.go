package types

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type Claims struct {
    UserID       uint
    jwt.StandardClaims
}

func (user User) GenerateJWT() (string, error) {
    claims := &Claims{
        UserID: user.Id,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signing_key := []byte(os.Getenv("JWT_SECRET"))
    token_string, err := token.SignedString(signing_key)

    return token_string, err
}
