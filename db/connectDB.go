package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/Naman-B-Parlecha/BullStash/config"
	_ "github.com/lib/pq"
)

func LoadPostgresDb() *sql.DB {
	var conf config.Config
	conf = *config.GetConfig()

	port, err := strconv.Atoi(conf.PORT)
	if err != nil {
		fmt.Println("Error parsing port")
		os.Exit(1)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.HOST, port, conf.USER, conf.PASSWORD, conf.DBNAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
