package postgres

import (
	"fmt"

	"github.com/Naman-B-Parlecha/BullStash/util"
)

func TestConnection(port int, dbname, host, user, password string) error {
	db, err := util.LoadPostgresDb(port, dbname, host, user, password)
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
