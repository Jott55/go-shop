package database

import (
	"context"
	"fmt"
	"jott55/go-shop/clog"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

type DatabaseInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type DatabaseLink struct {
	info *DatabaseInfo
	con  *pgx.Conn
}

type DatabaseResponse struct {
	str string
}

type Field struct {
	fieldsName []string
	fieldsAddr []any
}

func checkError(err error, msg ...any) bool {
	if err != nil {
		clog.Log(clog.ERROR, err, msg)
		return true
	}
	return false
}

func debug(msg ...any) {
	clog.Log(clog.DEBUG, msg...)
}

func Create() *DatabaseLink {
	return &DatabaseLink{}
}

func Configure(dl *DatabaseLink, dinfo DatabaseInfo) {
	dl.info = &dinfo
}

func Init(dl *DatabaseLink) error {
	db := dl.info

	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == "" || db.Database == "" {
		return fmt.Errorf("no database credentials")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Database)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return err
	}

	dl.con = conn

	return nil
}

func Exec(dl *DatabaseLink, sql string, args ...any) (DatabaseResponse, error) {
	comm, err := dl.con.Exec(context.Background(), sql, args...)
	return newDatabaseResponse(comm.String()), err
}

func Close(dl *DatabaseLink) error {
	return dl.con.Close(context.Background())
}

func QueryRow(dl *DatabaseLink, sql string, dest ...any) error {
	return dl.con.QueryRow(context.Background(), sql).Scan(dest...)
}

func Query(dl *DatabaseLink, sql string) (pgx.Rows, error) {
	return dl.con.Query(context.Background(), sql)
}

func newDatabaseResponse(str string) DatabaseResponse {
	return DatabaseResponse{str}
}

func CollectRows[T any](rows pgx.Rows) ([]T, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func GenericGet[T any](dl *DatabaseLink, table string, id int) (T, error) {
	var serial T

	res := getStructFields(&serial)

	items := strings.Join(res.fieldsName, ", ")
	sql_string := fmt.Sprintf("SELECT %v FROM %v WHERE id=%v", items, table, id)

	err := QueryRow(dl, sql_string, res.fieldsAddr...)

	if checkError(err) {
		return serial, err
	}

	debug(sql_string)

	debug(serial)

	return serial, nil
}

func GenericInsert[T any](dl *DatabaseLink, table string) {

}

func GenericGetWhere[T any](dl *DatabaseLink, table string, where string) []T {
	names := getStructNames[T]()

	cols := strings.Join(names, ", ")
	sql_string := fmt.Sprintf("SELECT %v FROM %v WHERE %v", cols, table, where)

	rows, err := Query(dl, sql_string)
	if checkError(err) {
		return nil
	}

	items, err := CollectRows[T](rows)

	if checkError(err) {
		return nil
	}

	return items
}

// Get addresses of the fields of the struct v
func getStructFields(v any) Field {
	structPointer := reflect.ValueOf(v) // struct pointer

	isPointer(structPointer.Kind())

	s := structPointer.Elem() // struct

	isStruct(s.Kind())

	length := s.NumField()

	var res Field

	t := s.Type() // struct type

	for i := range length {
		res.fieldsName = append(res.fieldsName, t.Field(i).Name)
		res.fieldsAddr = append(res.fieldsAddr, s.Field(i).Addr().Interface())
	}

	return res
}

func getStructNames[T any]() []string {

	t := reflect.TypeFor[T]() // type

	isStruct(t.Kind())

	length := t.NumField()

	var fieldsName []string

	for i := range length {
		fieldsName = append(fieldsName, t.Field(i).Name)
	}

	return fieldsName
}

func isPointer(v reflect.Kind) {
	if v != reflect.Pointer {
		panic(fmt.Sprintf("expected pointer to a struct, received %s", v))
	}
}

func isStruct(v reflect.Kind) {
	if v != reflect.Struct {
		panic(fmt.Sprintf("expected a struct, received %s", v))
	}
}
