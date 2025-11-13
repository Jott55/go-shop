package user

import (
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
)

type User struct {
	Id            int
	Name          string
	Email         string
	Password_hash string // plan on using argon2id
	Photo_url     string
}

func CreateTable(dl *database.DatabaseLink) {
	sql_table := `CREATE TABLE users (
		id bigint GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		email VARCHAR(50),
		password_hash VARCHAR(64),
		photo_url VARCHAR(255)
	)`

	table, err := database.Exec(dl, sql_table)

	if err != nil {
		clog.Log(clog.ERROR, err)
	}
	clog.Log(clog.DEBUG, table, sql_table)
}
