package database

import (
	"context"
	"fmt"
	"jott55/go-shop/clog"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IDatabase interface {
	Create()
	Configure(dinfo DatabaseInfo)
	Init() error
	Exec(sql string, args ...any) (DatabaseResponse, error)
	Close() error
	QueryRow(sql string, dest ...any) error
	Query(sql string) (pgx.Rows, error)
	Insert(table string, t any) DatabaseResponse
	DeleteById(table string, id int) error
	CreateTable(table string, fields string)
	CreateIndex(table string, columns []string)
	DropTable(table string)
	DeleteFromTableById(table string, id int)
	DeleteFromTableWhere(table string, condition string)
}

type DatabaseInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type DatabaseLink struct {
	info *DatabaseInfo
	con  *pgxpool.Pool
}

type Response interface {
	String() string
}

type DatabaseResponse struct {
	msg string
}

// TODO: change database reponse to be only interface or something like that
func (dr DatabaseResponse) String() string {
	return dr.msg
}

type FieldAddress struct {
	fieldName    []string
	fieldAddress []any
}

type FieldValues struct {
	fieldName  []string
	fieldValue []any
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

func Create() *DatabaseLink {
	return &DatabaseLink{}
}

func (dl *DatabaseLink) Configure(dinfo DatabaseInfo) {
	dl.info = &dinfo
}

func (dl *DatabaseLink) GetDBInfo() *DatabaseInfo {
	return dl.info
}

func (dl *DatabaseLink) Init() error {
	db := dl.info

	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == "" || db.Database == "" {
		return fmt.Errorf("no database credentials")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Database)
	// conn, err := pgx.Connect(context.Background(), url)

	pool, err := pgxpool.New(context.Background(), url)

	if err != nil {
		return err
	}
	err = pool.Ping(context.Background())

	if err != nil {
		return err
	}

	dl.con = pool

	return nil
}

func (dl *DatabaseLink) Exec(sql string, args ...any) (DatabaseResponse, error) {
	comm, err := dl.con.Exec(context.Background(), sql, args...)
	return DatabaseResponse{comm.String()}, err
}

func (dl *DatabaseLink) Close() error {
	dl.con.Close()
	return nil
}

func (dl *DatabaseLink) QueryRow(sql string, dest ...any) error {

	return dl.con.QueryRow(context.Background(), sql).Scan(dest...)
}

func (dl *DatabaseLink) Query(sql string) (pgx.Rows, error) {
	return dl.con.Query(context.Background(), sql)
}

func CollectRows[T any](rows pgx.Rows) ([]T, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func GenericGet[T any](dl *DatabaseLink, table string, id int) (T, error) {
	var serial T

	fa := getStructFieldsAddress(&serial)

	items := strings.Join(fa.fieldName, ", ")
	sql_string := fmt.Sprintf(`SELECT %v FROM %v WHERE id=%v`, items, table, id)

	err := dl.QueryRow(sql_string, fa.fieldAddress...)

	if checkError(err) {
		return serial, err
	}

	debug(sql_string)

	debug(serial)

	return serial, nil
}

// t = pointer to a struct
func (dl *DatabaseLink) Insert(table string, t any) DatabaseResponse {
	fv := getStructValues(t)

	cols := strings.Join(fv.fieldName, ", ")
	var values []string

	for _, val := range fv.fieldValue {
		values = append(values, fmt.Sprintf(`'%v'`, val))
	}

	valuesStr := strings.Join(values, ", ")

	sql_insert := fmt.Sprintf(`INSERT INTO %v (%v) VALUES (%v)`, table, cols, valuesStr)

	debug(sql_insert)

	tag, err := dl.Exec(sql_insert)

	checkError(err)

	return tag
}

func (dl *DatabaseLink) Update(table string, t any, where string) DatabaseResponse {
	fv := getStructValues(t)

	var values []string

	for i := range fv.fieldName {
		values = append(values, fmt.Sprintf(`%s='%v'`, fv.fieldName[i], fv.fieldValue[i]))
	}

	cols := strings.Join(values, ", ")

	sql_update := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, cols, where)

	debug(sql_update)
	tag, err := dl.Exec(sql_update)

	checkError(err)
	return tag
}

func GenericGetWhere[T any](dl *DatabaseLink, table string, where string) []T {
	names := getStructNames[T]()

	cols := strings.Join(names, ", ")
	sql_string := fmt.Sprintf(`SELECT %v FROM %v WHERE %v`, cols, table, where)

	debug(sql_string)

	rows, err := dl.Query(sql_string)
	if checkError(err) {
		return nil
	}

	items, err := CollectRows[T](rows)

	if checkError(err) {
		return nil
	}

	return items
}

// Get address of the fields of the struct v
func getStructFieldsAddress(v any) FieldAddress {
	var value = reflect.ValueOf(v)

	for value.Kind() != reflect.Struct {
		value = value.Elem()
	}

	t := value.Type()

	var res FieldAddress

	for i := range value.NumField() {
		res.fieldName = append(res.fieldName, t.Field(i).Name)
		res.fieldAddress = append(res.fieldAddress, value.Field(i).Addr().Interface())
	}
	return res
}

// Get the value of the fields of the struct *v
func getStructValues(v any) FieldValues {

	var value = reflect.ValueOf(v)

	for value.Kind() != reflect.Struct {
		value = value.Elem()
	}

	t := value.Type()

	var res FieldValues

	for i := range value.NumField() {
		res.fieldName = append(res.fieldName, t.Field(i).Name)
		res.fieldValue = append(res.fieldValue, value.Field(i).Interface())
	}
	return res
}

func getStructNames[T any]() []string {

	t := reflect.TypeFor[T]() // type

	// isStruct(t.Kind())

	length := t.NumField()

	var fieldsName []string

	for i := range length {
		fieldsName = append(fieldsName, t.Field(i).Name)
	}

	return fieldsName
}

// func isPointer(v reflect.Kind) bool {
// 	return v == reflect.Pointer
// }

// func isStruct(v reflect.Kind) bool {
// 	return v == reflect.Struct
// }

func (dl *DatabaseLink) DeleteById(table string, id int) error {
	sql_delete := fmt.Sprintf("DELETE FROM %v WHERE id=%v", table, id)

	_, err := dl.Exec(sql_delete)

	if checkError(err) {
		return err
	}
	return nil
}

func (dl *DatabaseLink) CreateTable(table string, fields string) {
	sql_table := fmt.Sprintf(`CREATE TABLE %v (%v)`, table, fields)

	dr, err := dl.Exec(sql_table)

	if checkError(err) {
		return
	}

	debug(dr, sql_table)
}

func (dl *DatabaseLink) CreateIndex(table string, columns []string) {
	cols := strings.Join(columns, ", ")
	sql_index := fmt.Sprintf(`CREATE INDEX %s_index ON %s (%v)`, table, table, cols)

	dr, err := dl.Exec(sql_index)

	if checkError(err) {
		return
	}

	debug(dr, sql_index)
}

func (dl *DatabaseLink) DropTable(table string) {
	sql_drop := fmt.Sprintf(`DROP TABLE %s`, table)

	dr, err := dl.Exec(sql_drop)

	if checkError(err) {
		return
	}

	debug(dr)

}

func (dl *DatabaseLink) DeleteFromTableById(table string, id int) {
	sql_delete := fmt.Sprintf(`DELETE FROM %s WHERE id=%d`, table, id)

	_, err := dl.Exec(sql_delete)
	checkError(err)
}

func (dl *DatabaseLink) DeleteFromTableWhere(table string, condition string) {
	sql_delete := fmt.Sprintf(`DELETE FROM %s WHERE %s`, table, condition)

	_, err := dl.Exec(sql_delete)
	checkError(err)
}
