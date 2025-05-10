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
		return fmt.Errorf("failed to restore database: %v", err)
	}

	successMsg := fmt.Sprintf("Database %s restored successfully from %s", dbname, input)
	util.CallWebHook(successMsg, false)
	util.SuccessColor.Println(successMsg)
	return nil
}
