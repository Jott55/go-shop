package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func Exec(dl *DatabaseLink, sql string) (pgconn.CommandTag, error) {
	return dl.con.Exec(context.Background(), sql)
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

func IsError(err error) {

}
