/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scheduling a backup job...")
		dbType, _ := cmd.Flags().GetString("dbtype")
		backupType, _ := cmd.Flags().GetString("backuptype")
		// outputDir, _ := cmd.Flags().GetString("output")
		cron, _ := cmd.Flags().GetString("cron")

		fmt.Println("Kindly put your values in env variables so that we can fetch from there")
		fmt.Print("Do you want to continue? (y/n): ")
		var answer string
		fmt.Scanln(&answer)

		if answer != "y" && answer != "Y" {
			util.CallWebHook("User cancelled the backup scheduling", true)
			fmt.Println("Operation cancelled")
			return
		}

		projectDir, err := os.Getwd()
		if err != nil {
			util.CallWebHook("Error getting current directory: "+err.Error(), true)
			fmt.Println("Error getting current directory:", err)
			return
		}

		fmt.Println("Current directory:", projectDir)

		for _, dir := range []string{"cron_job", "cron_logs"} {
			dirPath := filepath.Join(projectDir, dir)
			if err := os.MkdirAll(dirPath, 0755); err != nil && !os.IsExist(err) {
				util.CallWebHook("Error creating directory: "+err.Error(), true)
				fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
				return
			}
		}

		scriptPath := filepath.Join(projectDir, "cron_job", "backup_cron_job.sh")
		file, err := os.Create(scriptPath)
		if err != nil {
			util.CallWebHook("Error creating script file: "+err.Error(), true)
			fmt.Println("Error creating script file:", err)
			return
		}
		defer file.Close()

		scriptContent := fmt.Sprintf(`#!/bin/bash
cd "%s" || exit 1

echo "[$(date)] Starting BullStash backup..." >> "%s/cron_logs/backup.log"

BullStash backup --dbtype %s --backup-type %s --isCron true \
    >> "%s/cron_logs/backup.log" 2>&1

echo "[$(date)] Backup completed." >> "%s/cron_logs/backup.log"
`,
			projectDir, projectDir, dbType, backupType, projectDir, projectDir)

		if _, err := file.WriteString(scriptContent); err != nil {
			util.CallWebHook("Error writing to script file: "+err.Error(), true)
			fmt.Println("Error writing to file:", err)
			return
		}

		if err := os.Chmod(scriptPath, 0755); err != nil {
			util.CallWebHook("Error changing file permissions: "+err.Error(), true)
			fmt.Println("Error changing file permissions:", err)
			return
		}

		if cron == "" {
			cron = "* * * * *"
		}

		cronEntry := fmt.Sprintf("%s cd \"%s\" && \"%s\" >> \"%s/cron_logs/backup.log\" 2>&1",
			cron, projectDir, scriptPath, projectDir)

		addCronCmd := exec.Command("bash", "-c",
			fmt.Sprintf(`(crontab -l 2>/dev/null; echo "%s") | crontab -`, cronEntry))

		if output, err := addCronCmd.CombinedOutput(); err != nil {
			util.CallWebHook("Failed to add cron job: "+err.Error(), true)
			fmt.Printf("Failed to add cron job: %v\n", err)
			fmt.Printf("Command output: %s\n", string(output))
			return
		}

		util.CallWebHook("Backup job scheduled successfully", false)
		fmt.Println("Successfully scheduled backup with cron expression:", cron)
		fmt.Println("Script location:", scriptPath)
		fmt.Println("Logs will be written to:", filepath.Join(projectDir, "cron_logs/backup.log"))
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().String("cron", "", "Cron expression for scheduling the backup")
	scheduleCmd.Flags().String("dbtype", "postgres", "Type of database to backup")
	scheduleCmd.Flags().String("backuptype", "full", "Type of backup to perform (full/incremental)")
	scheduleCmd.Flags().String("output", "", "Output directory for the backup files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scheduleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
