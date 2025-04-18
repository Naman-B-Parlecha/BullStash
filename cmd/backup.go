/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/internal/config"
	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		output, _ := cmd.Flags().GetString("output")
		compress, _ := cmd.Flags().GetBool("compress")
		isCron, _ := cmd.Flags().GetBool("isCron")

		if dbtype != "postgres" {
			util.CallWebHook("Unsupported database type: "+dbtype, true)
			fmt.Fprintf(os.Stderr, "Unsupported database type: %s\n", dbtype)
			return
		}

		if isCron {
			postgresConfig := config.GetPostgresConfig()
			host = postgresConfig.HOST
			port, _ = strconv.Atoi(postgresConfig.PORT)
			user = postgresConfig.USER
			password = postgresConfig.PASSWORD
			dbname = postgresConfig.DBNAME
		}

		projectDir, err := os.Getwd()

		if err != nil {
			util.CallWebHook("Error getting current directory: "+err.Error(), true)
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}
		folderName := filepath.Join(projectDir, output)
		if err := os.MkdirAll(folderName, 0755); err != nil {
			util.CallWebHook("Error creating output directory: "+err.Error(), true)
			fmt.Printf("Failed to create output directory: %v\n", err)
			return
		}

		timestamp := time.Now().Format("20060102_150405")
		fileName := fmt.Sprintf("%s/%s_backup_%s.sql", output, dbname, timestamp)
		gzFileName := fileName + ".gz"

		sqlFile, err := os.Create(fileName)
		if err != nil {
			util.CallWebHook("Error creating backup file: "+err.Error(), true)
			fmt.Printf("Failed to create backup file: %v\n", err)
			return
		}
		defer sqlFile.Close()

		dumpCmd := exec.Command("pg_dump",
			"-h", host,
			"-p", strconv.Itoa(port),
			"-U", user,
			"-d", dbname)
		dumpCmd.Env = append(os.Environ(), "PGPASSWORD="+password)
		dumpCmd.Stdout = sqlFile

		if err := dumpCmd.Run(); err != nil {
			util.CallWebHook("pg_dump failed: "+err.Error(), true)
			fmt.Printf("pg_dump failed: %v\n", err)
			os.Remove(fileName)
			return
		}

		if compress {
			if err := util.CompressFile(fileName, gzFileName); err != nil {
				fmt.Printf("Compression failed: %v\n", err)
				return
			}
			if err := os.Remove(fileName); err != nil {
				util.CallWebHook("Error removing uncompressed file: "+err.Error(), true)
				fmt.Printf("Warning: could not remove uncompressed file: %v\n", err)
			}

			util.CallWebHook("Backup created successfully at: "+gzFileName, false)
			fmt.Printf("Backup successfully created at: %s\n", gzFileName)
			return
		}

		util.CallWebHook("Backup created successfully at: "+fileName, false)
		fmt.Printf("Backup successfully created at: %s\n", fileName)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().String("dbtype", "postgres", "Type of database")
	backupCmd.Flags().String("host", "", "Host of the database")
	backupCmd.Flags().Int("port", 0, "Port of the database")
	backupCmd.Flags().String("user", "", "User of the database")
	backupCmd.Flags().String("password", "", "Password of the database")
	backupCmd.Flags().String("dbname", "postgres", "Name of the database")
	backupCmd.Flags().String("backup-type", "full", "Type of backup such as full, incremental, differential")
	backupCmd.Flags().String("output", "backups", "Path you want to store the backup")
	backupCmd.Flags().Bool("compress", false, "Compress the backup file")
	backupCmd.Flags().String("storage", "local", "Storage type such as local, s3, gcs, goodle_drive")
	backupCmd.Flags().String("cloud-bucket", "", "Cloud bucket name")
	backupCmd.Flags().String("cloud-region", "asia-pacific-1", "Cloud region")
	backupCmd.Flags().Bool("isCron", false, "Is this a cron job")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
