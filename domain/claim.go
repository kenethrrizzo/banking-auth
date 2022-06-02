package domain

import "github.com/golang-jwt/jwt/v4"

type AccessTokenClaim struct {
	CustomerId string   `json:"customer_id"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Accounts   []string `json:"accounts"`
	jwt.StandardClaims
}

func (c AccessTokenClaim) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaim) IsValidCustomerId(customerId string) bool {
	return c.CustomerId == customerId
}