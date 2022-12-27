package sqlx

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dsnPattern = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"

func NewConnection(host, user, password, database string, port uint, sslMode bool) (*sqlx.DB, error) {
	ssl := "disable"
	if sslMode == true {
		ssl = "enable"
	}

	dsn := fmt.Sprintf(dsnPattern, host, port, user, password, database, ssl)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
