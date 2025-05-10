package mysql

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/Naman-B-Parlecha/BullStash/util"
)

func Restore(dbname, input, user, password string) error {
	fmt.Printf("Restoring database %s from %s\n", dbname, input)

	if _, err := os.Stat(input); os.IsNotExist(err) {
		errorMsg := fmt.Sprintf("Input file not found: %s", input)
		return fmt.Errorf("%s", errorMsg)
	}

	sqlFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("error opening SQL file: %w", err)
	}
	defer sqlFile.Close()

	restoreCmd := exec.Command("mysql",
		"-u", user,
		"-p"+password,
		dbname)

	restoreCmd.Stdin = sqlFile

	var stderr bytes.Buffer
	restoreCmd.Stderr = &stderr

	if err := restoreCmd.Run(); err != nil {
		errorMsg := stderr.String()
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		fullError := fmt.Sprintf("Failed to restore database: %s", errorMsg)
		return fmt.Errorf("%s", fullError)
	}

	successMsg := fmt.Sprintf("Database %s restored successfully from %s", dbname, input)
	util.CallWebHook(successMsg, false)
	util.SuccessColor.Println(successMsg)
	return nil
}
