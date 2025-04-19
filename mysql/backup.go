package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
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
		"-p"+password,
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
		return nil
	}
	util.CallWebHook("Backup created successfully at: "+fileName, false)
	fmt.Printf("Backup successfully created at: %s\n", fileName)
	return nil
}
