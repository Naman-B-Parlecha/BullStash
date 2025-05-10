package postgres

import (
	"github.com/Naman-B-Parlecha/BullStash/util"
)

func TestConnection(port int, dbname, host, user, password string) error {
	db, err := util.LoadPostgresDb(port, dbname, host, user, password)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
