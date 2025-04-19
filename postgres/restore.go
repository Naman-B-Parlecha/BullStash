package postgres

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/Naman-B-Parlecha/BullStash/util"
)

func Restore(dbname, input, host, user, password string, port int) error {
	fmt.Printf("Restoring database %s from %s\n", dbname, input)
	restoreCmd := exec.Command("psql",
		"-h", host,
		"-p", strconv.Itoa(port),
		"-U", user,
		"-d", dbname,
		"-f", input)

	restoreCmd.Env = append(restoreCmd.Env, "PGPASSWORD="+password)

	if err := restoreCmd.Run(); err != nil {
		util.CallWebHook("Failed to restore database: "+err.Error(), true)
		fmt.Printf("Failed to restore database: %v\n", err)
		return err
	}

	util.CallWebHook("Database restored successfully", false)
	fmt.Println("Database restored successfully")
	return nil
}
