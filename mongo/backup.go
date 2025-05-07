package mongo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/go-resty/resty/v2"
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

	client := resty.New()
	fileInfo, err := os.Stat(backupFolder)
	if err != nil {
		fmt.Printf("Error getting file size: %v\n", err)
	}
	fileSize := fileInfo.Size()

	_, err = client.R().SetBody(struct {
		DBType     string  `json:"dbtype"`
		BackupType string  `json:"backup_type"`
		Storage    string  `json:"storage"`
		Size       float64 `json:"size"`
	}{
		DBType:     "mongo",
		BackupType: "full",
		Storage:    "local",
		Size:       float64(fileSize),
	}).Post("http://localhost:8085/backups/size")

	fmt.Println("File size sent to monitoring service:", fileSize)
	if err != nil {
		fmt.Println("Error sending request:", err)
		util.CallWebHook("Error sending request: "+err.Error(), true)
	}

	successMsg := fmt.Sprintf("MongoDB backup created successfully at: %s", backupFolder)
	fmt.Println(successMsg)
	util.CallWebHook(successMsg, false)
	return nil
}
