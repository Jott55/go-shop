package cart

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"jott55/go-shop/services/user"
	"jott55/go-shop/types"
)

const Table = "cart"

const (
	const_id      = "id"
	const_user_id = "user_id"
)

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
		%s BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		%s BIGINT REFERENCES %s(id) ON DELETE CASCADE
	`, const_id, const_user_id, user.Table)

	database.CreateTable(dl, Table, field)
	database.CreateIndex(dl, Table, []string{"user_id"})
}

func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}

func Insert(dl *database.DatabaseLink, cart *types.CartNoId) {
	database.GenericInsert(dl, Table, cart)
}

func Delete(dl *database.DatabaseLink, id int) {
	database.DeleteById(dl, Table, id)
}

func Get(dl *database.DatabaseLink, id int) (types.Cart, error) {
	return database.GenericGet[types.Cart](dl, Table, id)
}

func GetIdByUserId(dl *database.DatabaseLink, user_id int) int {
	ar := database.GenericGetWhere[types.CartId](dl, Table, fmt.Sprintf("%s=%d", const_user_id, user_id))
	if len(ar) == 1 {
		return ar[0].Id
	}
	return 0
}

func GetAll(dl *database.DatabaseLink) []types.Cart {
	return database.GenericGetWhere[types.Cart](dl, Table, "true")
}
