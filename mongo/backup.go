package mongo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
)

func Backup(uri, dbName, outputDir string, isCompressed bool) error {
	if uri == "" {
		util.CallWebHook("MongoDB URI cannot be empty", true)
		return fmt.Errorf("MongoDB URI cannot be empty")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		util.CallWebHook("Error creating output directory: "+err.Error(), true)
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFolder := filepath.Join(outputDir, fmt.Sprintf("backup_%s_%s", dbName, timestamp))

	args := []string{
		"--uri=" + uri,
		"--out=" + backupFolder,
	}

	if dbName != "" {
		args = append(args, "--db="+dbName)
	}

	if isCompressed {
		args = append(args, "--gzip")
	}

	cmd := exec.Command("mongodump", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		errorMsg := fmt.Sprintf("mongodump failed: %v\nOutput: %s", err, string(output))
		fmt.Println(errorMsg)
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)
	}

	if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
		errorMsg := "backup folder was not created"
		fmt.Println(errorMsg)
		util.CallWebHook(errorMsg, true)
		return fmt.Errorf("%s", errorMsg)
	}

	successMsg := fmt.Sprintf("MongoDB backup created successfully at: %s", backupFolder)
	fmt.Println(successMsg)
	util.CallWebHook(successMsg, false)
	return nil
}
