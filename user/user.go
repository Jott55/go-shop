package user

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
)

const Table = "users"

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
		clog.Logger(clog.ERROR, 2, err, msg)
		return true
	}
	return false
}

func debug(msg ...any) {
	clog.Log(clog.DEBUG, msg...)
}

func CreateTable(dl *database.DatabaseLink) {
	field := `
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		email VARCHAR(50),
		password_hash VARCHAR(64),
		photo_url VARCHAR(255)
	`

	database.CreateTable(dl, Table, field)
}

func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}

func Get(dl *database.DatabaseLink, id int) (User, error) {
	return database.GenericGet[User](dl, Table, id)
}

func GetWhere(dl *database.DatabaseLink, id_min int, id_max int) []User {
	return database.GenericGetWhere[User](dl, Table, fmt.Sprintf("id BETWEEN %v AND %v", id_min, id_max))
}

func Insert(dl *database.DatabaseLink, user *UserInsert) database.DatabaseResponse {
	return database.GenericInsert(dl, Table, user)
}

func Delete(dl *database.DatabaseLink, id int) error {
	return database.DeleteById(dl, Table, id)
}
