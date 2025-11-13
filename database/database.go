package database

import (
	"context"
	"fmt"
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
