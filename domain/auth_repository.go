package domain

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type AuthRepository interface {
	FindBy(string, string) (*Login, error)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, error) {
	var login Login
	sqlVerify := `
		select u.Username, u.CustomerId, u.Role, group_concat(a.Id) as AccountNumbers 
		from Users u
		left join Accounts a on a.Id = u.CustomerId
		where u.Username = ? and u.Password = ?
		group by a.CustomerId
	`
	err := d.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid credentials")
		}
		log.Error("Error while verifying login request from database: ", err.Error())
		return nil, errors.New("Unexpected database error")
	}
	return &login, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
