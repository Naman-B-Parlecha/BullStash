/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/config"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var postgres_config config.Config
		postgres_config = *config.GetPostgresConfig()

		dbtype, _ := cmd.Flags().GetString("dbtype")
		output, _ := cmd.Flags().GetString("output")

		if err := os.MkdirAll(output, 0755); err != nil {
			fmt.Errorf("failed to create output directory: %v", err.Error())
			return
		}

		backupFile := fmt.Sprintf("%s/%sbackup_%s.sql", output, postgres_config.DBNAME, time.Now().UTC().Local())

		fmt.Println("Backup of " + dbtype + " database is stored in " + backupFile)

		command := exec.Command("pg_dump",
			"-h", postgres_config.HOST,
			"-p", postgres_config.PORT,
			"-U", postgres_config.USER,
			"-d", postgres_config.DBNAME,
			"-f", backupFile)

		command.Env = append(command.Env, fmt.Sprintf("PGPASSWORD=%s", postgres_config.PASSWORD))

		outputDetails, err := command.CombinedOutput()

		if err != nil {
			fmt.Println("Something went wrong in the backup :", err.Error())
			return
		}

		fmt.Println(string(outputDetails))

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().String("dbtype", "postgres", "Type of database")
	backupCmd.Flags().String("output", "backup.sql", "Path you want to store the backup")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
