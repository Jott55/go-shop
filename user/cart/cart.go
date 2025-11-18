package cart

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"jott55/go-shop/user"
)

const Table = "users_cart"

type Cart struct {
	Id      int64
	User_id int
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
	field := fmt.Sprintf(`
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		user_id BIGINT REFERENCES %s(id) ON DELETE CASCADE
	`, user.Table)

	database.CreateTable(dl, Table, field)
	database.CreateIndex(dl, Table, []string{"user_id"})
}

func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}

func Insert(dl *database.DatabaseLink, cart *Cart) {
	database.GenericInsert(dl, Table, cart)
}

func Delete(dl *database.DatabaseLink, id int) {
	database.DeleteById(dl, Table, id)
}
