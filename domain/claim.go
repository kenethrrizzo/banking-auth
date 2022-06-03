package domain

import (
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type AccessTokenClaim struct {
	CustomerId string   `json:"CustomerId"`
	Username   string   `json:"Username"`
	Role       string   `json:"Role"`
	Accounts   []string `json:"Accounts"`
	jwt.StandardClaims
}

func (c AccessTokenClaim) IsUserRole() bool {
	return c.Role == "user"
}

func (c AccessTokenClaim) IsValidCustomerId(customerId string) bool {
	return c.CustomerId == customerId
}

func (c AccessTokenClaim) IsValidAccountId(accountId string) bool {
	log.Info("AccountID: ", accountId)
	log.Info("Accounts in token: ", c.Accounts)
	if accountId == "" {
		return true
	}
	for _, a := range c.Accounts {
		if a == accountId {
			return true
		}
	}
	return false
}

func (c AccessTokenClaim) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	return c.CustomerId == urlParams["customer-id"] && c.IsValidAccountId(urlParams["account-id"])
}
