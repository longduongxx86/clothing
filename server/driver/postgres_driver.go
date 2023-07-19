package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	SQL *sql.DB
}

func ConnectPsql(host, port, user, password, dbname string, maxIdleConns, maxOpenConns int) *PostgresDB {

	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, _ := sql.Open("postgres", connectString)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	return &PostgresDB{
		SQL: db,
	}
}
