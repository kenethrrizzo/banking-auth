package domain

import "database/sql"

type Login struct {
	Username   string `db:"Username"`
	CustomerId sql.NullString `db:"CustomerId"`
	Accounts sql.NullString `db:"Accounts"`
	//TODO
}