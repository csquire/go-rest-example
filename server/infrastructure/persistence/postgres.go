package persistence

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDb(host, port, user, password, dbname, sslmode string) (*sqlx.DB, error) {

	if user == "" || host == "" {
		return nil, errors.New("user and host are required arguments")
	}
	connStr := fmt.Sprintf("user=%s host=%s dbname=%s connect_timeout=10", user, host, dbname)

	if password != "" {
		connStr = connStr + fmt.Sprintf(" password=%s", password)
	}

	if port != "" {
		connStr = connStr + fmt.Sprintf(" port=%s", port)
	}

	if sslmode != "" {
		connStr = connStr + fmt.Sprintf(" sslmode=%s", sslmode)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	return db, nil
}
