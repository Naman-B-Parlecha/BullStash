package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/go-resty/resty/v2"
)

func Backup(output, dbname, user, password string, compress bool) error {

	projectDir, err := os.Getwd()

	if err != nil {
		util.CallWebHook("Error getting current directory: "+err.Error(), true)
		fmt.Printf("Error getting current directory: %v\n", err)
		return err
	}
	folderName := filepath.Join(projectDir, output)
	if err := os.MkdirAll(folderName, 0755); err != nil {
		util.CallWebHook("Error creating output directory: "+err.Error(), true)
		fmt.Printf("Failed to create output directory: %v\n", err)
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s_backup_%s.sql", output, dbname, timestamp)
	gzFileName := fileName + ".gz"

	sqlFile, err := os.Create(fileName)
	if err != nil {
		util.CallWebHook("Error creating backup file: "+err.Error(), true)
		fmt.Printf("Failed to create backup file: %v\n", err)
		return err
	}
	defer sqlFile.Close()

	dumpCmd := exec.Command("mysqldump",
		"-u", user,
		"--password="+password,
		dbname)

	dumpCmd.Stdout = sqlFile

	if err := dumpCmd.Run(); err != nil {
		util.CallWebHook("my_sql_dump failed: "+err.Error(), true)
		fmt.Printf("my_sql_dump failed: %v\n", err)
		os.Remove(fileName)
		return err
	}

	if compress {
		if err := util.CompressFile(fileName, gzFileName); err != nil {
			fmt.Printf("Compression failed: %v\n", err)
			return err
		}
		if err := os.Remove(fileName); err != nil {
			util.CallWebHook("Error removing uncompressed file: "+err.Error(), true)
			fmt.Printf("Warning: could not remove uncompressed file: %v\n", err)
		}

		util.CallWebHook("Backup created successfully at: "+gzFileName, false)
		fmt.Printf("Backup successfully created at: %s\n", gzFileName)

		client := resty.New()
		fileInfo, err := os.Stat(gzFileName)
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
			DBType:     "mysql",
			BackupType: "full",
			Storage:    "local",
			Size:       float64(fileSize),
		}).Post("http://localhost:8085/backups/size")

		fmt.Println("File size sent to monitoring service:", fileSize)
		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}

		return nil
	}

	client := resty.New()
	fileInfo, err := os.Stat(fileName)
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
		DBType:     "mysql",
		BackupType: "full",
		Storage:    "local",
		Size:       float64(fileSize),
	}).Post("http://localhost:8085/backups/size")

	fmt.Println("File size sent to monitoring service:", fileSize)
	if err != nil {
		fmt.Println("Error sending request:", err)
		util.CallWebHook("Error sending request: "+err.Error(), true)
	}

	util.CallWebHook("Backup created successfully at: "+fileName, false)
	fmt.Printf("Backup successfully created at: %s\n", fileName)
	return nil
}
