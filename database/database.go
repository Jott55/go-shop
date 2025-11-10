package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// aways remember to Close
func Init(db Database) (*pgx.Conn, error) {

	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == "" || db.Database == "" {
		return nil, fmt.Errorf("no database credentials")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Database)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Exec(conn *pgx.Conn, sql string) (pgconn.CommandTag, error) {
	return conn.Exec(context.Background(), sql)
}

func Close(conn *pgx.Conn) error {
	return conn.Close(context.Background())
}

func QueryRow(conn *pgx.Conn, sql string, dest ...any) error {
	return conn.QueryRow(context.Background(), sql).Scan(dest...)
}

func Query(conn *pgx.Conn, sql string) (pgx.Rows, error) {
	return conn.Query(context.Background(), sql)
}
