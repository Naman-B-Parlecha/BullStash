package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnection(dbname, user, password string) error {
	db, err := loadMySQLDb(dbname, user, password)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func loadMySQLDb(dbname string, user string, password string) (*sql.DB, error) {
	mysqlInfo := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	db, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
