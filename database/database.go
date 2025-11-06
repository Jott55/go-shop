package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database struct {
	Name     string
	Password string
	Host     string
	Port     string
	Database string
}

// aways remember to Close
func Init(db Database) *pgx.Conn {

	if db.Name == "" || db.Password == "" || db.Host == "" || db.Port == "" || db.Database == "" {
		fmt.Println("ERROR: no database credentials")
		os.Exit(1)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.Name, db.Password, db.Host, db.Port, db.Database)

	conn := tryConnect(url)

	return conn
}

func Exec(conn *pgx.Conn, sql string) (pgconn.CommandTag, error) {
	return conn.Exec(context.Background(), sql)
}

func Close(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func QueryRow(conn *pgx.Conn, sql string, dest ...any) error {
	return conn.QueryRow(context.Background(), sql).Scan(dest...)
}

func Query(conn *pgx.Conn, sql string) (pgx.Rows, error) {
	return conn.Query(context.Background(), sql)
}

func tryConnect(url string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), url)
	handleError(err)
	if err == nil {
		return conn
	}
	fmt.Println("\n\nEnter 'x' to retry")
	var ch rune
	for {
		fmt.Scanf("%c", &ch)

		switch ch {
		case 'x':
			return conn
		default:
			os.Exit(1)
		}
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("\nERROR: %v\n", err)
	}
}
