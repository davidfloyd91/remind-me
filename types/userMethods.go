package types

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (user User) GenerateJwt(isAdmin bool) (string, error) {
	claims := &Claims{
		UserId: user.Id,
		Admin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}
