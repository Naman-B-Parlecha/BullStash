package mysql

import (
	"database/sql"
	"fmt"

	"github.com/Naman-B-Parlecha/BullStash/util"

	_ "github.com/go-sql-driver/mysql"
)

func TestConnection(dbname, user, password string) error {
	db, err := loadMySQLDb(dbname, user, password)
	if err != nil {
		util.CallWebHook("Failed to connect to database: "+err.Error(), true)
		fmt.Printf("Failed to connect to database: %v\n", err)
		return err
	}
	util.CallWebHook("Database is connected successfully", false)
	fmt.Printf("Database is connected successfully\n")
	defer db.Close()
	return nil
}

func loadMySQLDb(dbname string, user string, password string) (*sql.DB, error) {
	mysqlInfo := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	db, err := sql.Open("mysql", mysqlInfo)
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
