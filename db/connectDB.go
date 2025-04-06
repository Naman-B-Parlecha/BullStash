package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func LoadPostgresDb(port int, dbname string, host string, user string, password string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
