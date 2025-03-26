/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		// backupType, _ := cmd.Flags().GetString("backup-type")
		output, _ := cmd.Flags().GetString("output")
		compress, _ := cmd.Flags().GetBool("compress")
		// storage, _ := cmd.Flags().GetString("storage")
		// cloudBucket, _ := cmd.Flags().GetString("cloud-bucket")
		// cloudRegion, _ := cmd.Flags().GetString("cloud-region")

		if dbtype != "postgres" {
			fmt.Printf("Unsupported database type: %s\n", dbtype)
			return
		}

		if err := os.MkdirAll(output, 0755); err != nil {
			fmt.Printf("Failed to create output directory: %v\n", err)
			return
		}

		timestamp := time.Now().Format("20060102_150405")
		fileName := fmt.Sprintf("%s/%s_backup_%s.sql", output, dbname, timestamp)
		gzFileName := fileName + ".gz"

		sqlFile, err := os.Create(fileName)
		if err != nil {
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
				fmt.Printf("Warning: could not remove uncompressed file: %v\n", err)
			}

			fmt.Printf("Backup successfully created at: %s\n", gzFileName)
			return
		}

		fmt.Printf("Backup successfully created at: %s\n", fileName)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().String("dbtype", "postgres", "Type of database")
	backupCmd.Flags().String("host", "localhost", "Host of the database")
	backupCmd.Flags().Int("port", 5432, "Port of the database")
	backupCmd.Flags().String("user", "postgres", "User of the database")
	backupCmd.Flags().String("password", "password", "Password of the database")
	backupCmd.Flags().String("dbname", "postgres", "Name of the database")
	backupCmd.Flags().String("backup-type", "full", "Type of backup such as full, incremental, differential")
	backupCmd.Flags().String("output", "backup", "Path you want to store the backup")
	backupCmd.Flags().Bool("compress", false, "Compress the backup file")
	backupCmd.Flags().String("storage", "local", "Storage type such as local, s3, gcs, goodle_drive")
	backupCmd.Flags().String("cloud-bucket", "", "Cloud bucket name")
	backupCmd.Flags().String("cloud-region", "asia-pacific-1", "Cloud region")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
