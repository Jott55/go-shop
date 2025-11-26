package serverio

import (
	"encoding/json"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func noDb(dl *database.DatabaseLink) bool {
	if dl == nil {
		clog.Log(clog.ERROR, "no db connection, returning")
		return true
	}
	return false
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

func SendJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(v)
	if checkError(err) {
		return
	}
	w.Write(content)
}
func GetId(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

// Accept any request struct
func GetStructFromRequestBody[T any](r *http.Request) (T, error) {

	var struc T

	err := json.NewDecoder(r.Body).Decode(&struc)

	if checkError(err) {
		return struc, err
	}

	debug(struc)

	return struc, nil
}
