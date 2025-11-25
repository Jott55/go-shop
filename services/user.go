package services

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/types"
)

type UserService struct {
	table string
	dl    *database.DatabaseLink
}

func (u *UserService) Init(dl *database.DatabaseLink, table_name string) {
	u.dl = dl
	u.table = table_name
}

func (u *UserService) Get(id int) (types.User, error) {
	return database.GenericGet[types.User](u.dl, u.table, id)
}

func (u *UserService) GetWhere(id_min int, id_max int) []types.User {
	return database.GenericGetWhere[types.User](u.dl, u.table, fmt.Sprintf("id BETWEEN %v AND %v", id_min, id_max))
}

func (u *UserService) Insert(user *types.UserNoId) database.DatabaseResponse {
	return u.dl.Insert(u.table, user)
}

func (u *UserService) Drop() {
	u.dl.DropTable(u.table)
}

func (u *UserService) Create() {
	field := `
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		email VARCHAR(50),
		password VARCHAR(64),
		photo_url VARCHAR(255)
	`

	u.dl.CreateTable(u.table, field)
}

func (u *UserService) Delete(id int) error {
	return u.dl.DeleteById(u.table, id)
}
