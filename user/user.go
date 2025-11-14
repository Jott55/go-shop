package user

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
)

const table = "users"

type User struct {
	Id            int
	Name          string
	Email         string
	Password_hash string // plan on using argon2id
	Photo_url     string
}

type UserInsert struct {
	Name          string
	Email         string
	Password_hash string
	Photo_url     string
}

func checkError(err error, msg ...any) bool {
	if err != nil {
		clog.Log(clog.ERROR, msg...)
		return true
	}
	return false
}

func debug(msg ...any) {
	clog.Log(clog.DEBUG, msg...)
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

	if checkError(err) {
		return
	}

	debug(table, sql_table)
}

func Get(dl *database.DatabaseLink, id int) (User, error) {
	return database.GenericGet[User](dl, table, id)
}

func GetWhere(dl *database.DatabaseLink, id_min int, id_max int) []User {
	return database.GenericGetWhere[User](dl, table, fmt.Sprintf("id BETWEEN %v AND %v", id_min, id_max))
}

func Insert(dl *database.DatabaseLink, user *UserInsert) database.DatabaseResponse {
	return database.GenericInsert(dl, table, user)
}
