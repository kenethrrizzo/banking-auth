package domain

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

const HMAC_SAMPLE_SECRET = "verySecretString"

type Login struct {
	Username   string         `db:"Username"`
	CustomerId sql.NullString `db:"CustomerId"`
	Accounts   sql.NullString `db:"AccountNumbers"`
	Role       string         `db:"Role"`
}

func (l Login) GenerateToken() (*string, error) {
	var claims jwt.MapClaims
	if l.Accounts.Valid && l.CustomerId.Valid {
		claims = l.claimsForUser()
	} else {
		claims = l.claimsForAdmin()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsStr, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Error("Failed while signing token:", err.Error())
		return nil, errors.New("Cannot generate token")
	}
	return &signedTokenAsStr, nil
}

func (l Login) claimsForUser() jwt.MapClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return jwt.MapClaims{
		"CustomerId": l.CustomerId.String,
		"Role":       l.Role,
		"Username":   l.Username,
		"Accounts":   accounts,
		"Exp":        time.Now().Add(time.Hour).Unix(),
	}
}

func (l Login) claimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"Role":     l.Role,
		"Username": l.Username,
		"Exp":      time.Now().Add(time.Hour).Unix(),
	}
}

